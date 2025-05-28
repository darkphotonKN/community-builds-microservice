package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/darkphotonKN/community-builds-microservice/example-service/config"
	"github.com/darkphotonKN/community-builds-microservice/item-service/internal/item"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	serviceName = "item-service"
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

	repo := item.NewRepository(db)
	service := item.NewService(repo)
	handler := item.NewHandler(service)

	go service.InitCrawling(db)

	pb.RegisterItemServiceServer(grpcServer, handler)

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

// package main

// import (
// 	// "context"
// 	"fmt"

// 	"github.com/darkphotonKN/community-builds-microservice/items-service/config"
// 	"github.com/darkphotonKN/community-builds-microservice/item-service/internal/item"

// 	// "log"
// 	"net/http"
// 	"os"

// 	// "github.com/jackc/pgx/v5"
// 	"github.com/joho/godotenv"
// )

// func main() {

// 	// env setup
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Println("No .env file found, using system environment variables")
// 	}
// 	db := config.InitDB()
// 	defer db.Close()

// 	itemRepository := item.NewItemRepository(db)
// 	itemService := item.NewItemService(itemRepository)
// 	itemService.CreateGrpcServer()

// 	// conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost:5857/item_service_db")
// 	// if err != nil {
// 	// 	log.Fatal("Unable to connect to database:", err)
// 	// }
// 	// defer conn.Close(context.Background())

// 	// var greeting string
// 	// err = conn.QueryRow(context.Background(), "SELECT 'Hello, PostgreSQL!'").Scan(&greeting)
// 	// if err != nil {
// 	// 	log.Fatal("Query failed:", err)
// 	// }

// 	defaultDevPort := ":700"
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultDevPort
// 	}

// 	err := http.ListenAndServe(defaultDevPort, nil)
// 	if err != nil {
// 		// panic(err)
// 	}
// }
