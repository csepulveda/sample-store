package main

import (
	"context"
	"log"
	"os"

	"orders-service/internal/handlers"
	"orders-service/internal/publisher"
	"orders-service/internal/repository"
	"orders-service/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	tempoEndpoint := getEnv("TEMPO_ENDPOINT", "tempo:4318")
	tp := tracing.InitTracer("orders-service", tempoEndpoint)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
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

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	tableName := getEnv("ORDERS_TABLE", "orders")
	orderRepo := repository.NewDynamoOrderRepository(dynamoClient, tableName)

	snsTopicArn := getEnv("ORDERS_TOPIC_ARN", "")
	if snsTopicArn == "" {
		log.Fatal("ORDERS_TOPIC_ARN environment variable is required")
	}

	snsClient := sns.NewFromConfig(cfg)
	orderPublisher := publisher.NewSnsOrderPublisher(snsClient, snsTopicArn)

	app := fiber.New()
	app.Use(otelfiber.Middleware())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})
	api := app.Group("/api/orders")
	api.Post("/", handlers.CreateOrderHandler(orderRepo, orderPublisher))
	api.Get("/", handlers.ListOrdersHandler(orderRepo))
	api.Get("/:id", handlers.ListOrdersHandler(orderRepo))
	api.Patch("/:id", handlers.PatchOrderHandler(orderRepo, orderPublisher))
	api.Delete("/:id", handlers.DeleteOrderHandler(orderRepo, orderPublisher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
