package main

import (
	"context"
	"log"
	"os"

	"orders-service/internal/handlers"
	"orders-service/internal/repository"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	tableName := os.Getenv("ORDERS_TABLE")
	if tableName == "" {
		log.Fatal("ORDERS_TABLE env var is required")
	}

	orderRepo := repository.NewDynamoOrderRepository(dynamoClient, tableName)

	api := app.Group("/api/orders")
	api.Post("/", handlers.CreateOrderHandler(orderRepo))
	api.Get("/", handlers.ListOrdersHandler(orderRepo))
	api.Get("/:id", handlers.ListOrdersHandler(orderRepo))
	api.Patch("/:id", handlers.PatchOrderHandler(orderRepo))
	api.Delete("/:id", handlers.DeleteOrderHandler(orderRepo))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
