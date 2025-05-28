package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/config"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/validation"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery/consul"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName            = "api-gateway"
	httpAddr               = commonhelpers.GetEnvString("PORT", "7001")
	exampleServiceGrpcAddr = commonhelpers.GetEnvString("GRPC_EXAMPLE_ADDR", "7010")
	itemServiceGrpcAddr    = commonhelpers.GetEnvString("GRPC_ITEM_SERVICE_ADDR", "7001")
	consulAddr             = commonhelpers.GetEnvString("CONSUL_ADDR", "192.168.0.207:8510")
)

/**
* Main entry point to entire application.
* NOTE: Keep code here as clean and little as possible.
**/
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
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+httpAddr); err != nil {
		fmt.Printf("\nError when registering service:\n\n%s\n\n", err)
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

	// --- router setup ---
	router := config.SetupRouter(registry, db)

	// -- custom validators --
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterValidators(v)
	}

	// -- start server --
	if err := router.Run(fmt.Sprintf(":%s", httpAddr)); err != nil {
		log.Fatal("Failed to start server")
	}
}
