package main

import (
	"context"
	"log"
	"os"
	"products-worker/internal/processor"
	"products-worker/internal/repository"
	"products-worker/internal/sqs"
	"products-worker/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	tempoEndpoint := getEnv("TEMPO_ENDPOINT", "tempo:4318")
	tp := tracing.InitTracer("products-worker", tempoEndpoint)
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

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	queueURL := getEnv("SQS_QUEUE_URL", "")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}

	tableName := getEnv("DYNAMODB_TABLE", "products")

	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo := repository.NewDynamoProductRepository(dynamoClient, tableName)

	sqsClient := awssqs.NewFromConfig(cfg)
	handler := processor.NewOrderHandler(repo)

	log.Println("Worker started. Listening for messages...")
	sqs.ListenAndProcess(context.Background(), sqsClient, queueURL, handler)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
