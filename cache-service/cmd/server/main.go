package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/darkphotonKN/community-builds-microservice/cache-service/internal/cache"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
	"github.com/darkphotonKN/community-builds-microservice/common/broker"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// grpc
	serviceName = "cache"
	grpcAddr    = commonhelpers.GetEnvString("GRPC_EXAMPLE_ADDR", "7005")
	consulAddr  = commonhelpers.GetEnvString("CONSUL_ADDR", "localhost:8510")

	// rabbit mq
	amqpUser     = commonhelpers.GetEnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonhelpers.GetEnvString("RABBITMQ_PASS", "guest")
	amqpHost     = commonhelpers.GetEnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = commonhelpers.GetEnvString("RABBITMQ_PORT", "5672")
)

func main() {
	fmt.Println("Hello, World!")
	// --- database setup ---

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

	// --- server initialization ---
	grpcServer := grpc.NewServer()

	// 啟用 gRPC reflection
	reflection.Register(grpcServer)

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

	broker.DeclareExchange(ch, commonconstants.ItemCreatedItemEvent, "fanout")
	defer func() {
		close()
		ch.Close()
	}()

	// Note: These lines reference 'item' package and 'db' which are not imported/defined
	// You'll need to import the item package and uncomment the db initialization
	// repo := item.NewRepository(db)
	// service := item.NewService(repo, ch)
	// handler := item.NewHandler(service)
	// go service.InitCrawling(db)

	// pb.RegisterItemServiceServer(grpcServer, handler)

	service := cache.NewService()
	handler := cache.NewHandler(service)

	pb.RegisterCacheServiceServer(grpcServer, handler)

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}
}
