package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/darkphotonKN/community-builds-microservice/common/broker"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/payment-service/config"
	"github.com/darkphotonKN/community-builds-microservice/payment-service/internal/grpc-gateway/auth"
	"github.com/darkphotonKN/community-builds-microservice/payment-service/internal/payment"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/stripe/stripe-go/v82"
)

var (
	serviceName = "payment-service"
	grpcAddr    = commonhelpers.GetEnvString("GRPC_PAYMENT_ADDR", "7005")
	consulAddr  = commonhelpers.GetEnvString("CONSUL_ADDR", "localhost:8510")
	// rabbit mq
	amqpUser     = commonhelpers.GetEnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonhelpers.GetEnvString("RABBITMQ_PASS", "guest")
	amqpHost     = commonhelpers.GetEnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = commonhelpers.GetEnvString("RABBITMQ_PORT", "5672")
)

func main() {

	// --- database setup ---

	db := config.InitDB()
	defer db.Close()

	// setup stripe
	stripe.Key = commonhelpers.GetEnvString("STRIPE_SECRET_KEY", "")

	// --- service discovery setup ---

	// -- consul client --
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		log.Fatal("Failed to create Consul registry")
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	// -- discovery --
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+grpcAddr); err != nil {
		log.Printf("\nError when registering service:\n\n%s\n\n", err)
		panic(err)
	}

	// -- health check --
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Health check failed.")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// --- server initialization ---
	grpcServer := grpc.NewServer()

	// create a network listener to this service
	listener, err := net.Listen("tcp", "localhost:"+grpcAddr)

	if err != nil {
		log.Fatalf(
			"Failed to listen at port: %s\nError: %s\n", grpcAddr, err,
		)
	}

	defer listener.Close()

	// --- message broker - rabbit mq ---
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	broker.DeclareExchange(ch, commonconstants.MemberSignedUpEvent, "fanout")

	broker.DeclareExchange(ch, commonconstants.RatingCreatedEvent, "fanout")
	defer func() {
		close()
		ch.Close()
	}()

	authClient := auth.NewClient(registry)
	authHandler := auth.NewHandler(authClient)
	stripeProcessor := payment.NewStripeProcessor()
	paymentService := payment.NewService(authHandler, stripeProcessor)
	paymentHandler := payment.NewHandler(paymentService)

	pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}

}
