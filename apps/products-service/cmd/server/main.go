package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"

	"products-service/internal/handlers"
	"products-service/internal/repository"
	"products-service/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	tempoEndpoint := getEnv("TEMPO_ENDPOINT", "tempo:4318")
	tp := tracing.InitTracer("products-service", tempoEndpoint)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	productsTable := getEnv("PRODUCTS_TABLE", "products")

	awsEndpoint := getEnv("AWS_ENDPOINT", "")
	awsRegion := getEnv("AWS_REGION", "us-west-2")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Cargar configuración AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	productRepo := repository.NewDynamoProductRepository(dynamoClient, productsTable)

	app := fiber.New()

	app.Use(otelfiber.Middleware())

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
