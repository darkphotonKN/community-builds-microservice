package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	"github.com/darkphotonKN/community-builds-microservice/common/broker"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	// "github.com/darkphotonKN/community-builds-microservice/notification-service/config"
	"github.com/darkphotonKN/community-builds-microservice/notification-service/internal/example"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	// grpc
	serviceName = "notifications"
	grpcAddr    = commonhelpers.GetEnvString("PORT", "7009")
	consulAddr  = commonhelpers.GetEnvString("CONSUL_ADDR", "localhost:8510")

	// rabbit mq
	amqpUser     = commonhelpers.GetEnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonhelpers.GetEnvString("RABBITMQ_PASS", "guest")
	amqpHost     = commonhelpers.GetEnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = commonhelpers.GetEnvString("RABBITMQ_PORT", "5672")
)

func main() {
	// --- database setup ---

	// TODO: removed DB for now, can re-add back later if needed
	// db := config.InitDB()
	// defer db.Close()

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

	// --- grpc ---
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

	defer func() {
		close()
		ch.Close()
	}()

	// TODO: removed DB for now, can re-add back later if needed
	// repo := example.NewRepository(db)
	service := example.NewService(nil, ch)
	handler := example.NewHandler(service)
	consumer := example.NewConsumer(service, ch)
	// start goroutine and listen to events from message broker
	consumer.Listen()

	pb.RegisterExampleServiceServer(grpcServer, handler)

	log.Printf("grpc Notification Server started on PORT: %s\n", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}
}

