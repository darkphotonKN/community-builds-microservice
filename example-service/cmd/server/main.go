package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/darkphotonKN/community-builds-microservice/example-service/config"
	"github.com/darkphotonKN/community-builds-microservice/example-service/internal/example"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	serviceName = "examples"
	grpcAddr    = commonhelpers.GetEnvString("GRPC_EXAMPLE_ADDR", "7010")
	consulAddr  = commonhelpers.GetEnvString("CONSUL_ADDR", "localhost:8510")
)

func main() {

	// --- database setup ---

	db := config.InitDB()
	defer db.Close()

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

	repo := example.NewRepository(db)
	service := example.NewService(repo)
	handler := example.NewHandler(service)

	pb.RegisterExampleServiceServer(grpcServer, handler)

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}

	/*
	   // service setup
	   repo := order.NewRepository(db)
	   service := order.NewService(repo, ch)

	   // start grpc server
	   handler := order.NewGrpcHandler(service)

	   // create server
	   pb.RegisterOrderServiceServer(grpcServer, handler)

	   log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)
	   // start serving requests

	   	if err := grpcServer.Serve(l); err != nil {
	   		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	   	}
	*/
}
