package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"products-service/internal/handlers"
	"products-service/internal/repository"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	awsRegion := getEnv("AWS_REGION", "us-west-2")
	productsTable := getEnv("PRODUCTS_TABLE", "products")

	// Cargar configuración AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	productRepo := repository.NewDynamoProductRepository(dynamoClient, productsTable)

	app := fiber.New()

	// Middleware Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	// Healthcheck
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Product Routes
	api := app.Group("/api/products")
	api.Post("/", handlers.CreateProductHandler(productRepo))
	api.Get("/:id?", handlers.ListProductsHandler(productRepo))
	api.Put("/:id", handlers.UpdateProductHandler(productRepo))
	api.Patch("/:id", handlers.PatchProductHandler(productRepo))
	api.Delete("/:id", handlers.DeleteProductHandler(productRepo))

	port := getEnv("PORT", "8080")
	log.Printf("Starting Products Service on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
