package repository

import (
	"context"
	"errors"
	"products-service/internal/domain"
	"products-service/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoProductRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoProductRepository(client *dynamodb.Client, tableName string) *DynamoProductRepository {
	return &DynamoProductRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoProductRepository) Create(ctx context.Context, product *domain.Product) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#Create")
	defer span.End()
	span.SetAttributes(
		tracing.StringAttribute("productId", product.ID),
	)

	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#GetAll")
	defer span.End()

	output, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	err = attributevalue.UnmarshalListOfMaps(output.Items, &products)
	return products, err
}

func (r *DynamoProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#GetByID")
	defer span.End()
	span.SetAttributes(
		tracing.StringAttribute("productId", id),
	)

	output, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if output.Item == nil || len(output.Item) == 0 {
		return nil, errors.New("product not found")
	}

	var product domain.Product
	err = attributevalue.UnmarshalMap(output.Item, &product)
	return &product, err
}

func (r *DynamoProductRepository) Update(ctx context.Context, product *domain.Product) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#Update")
	defer span.End()

	return r.Create(ctx, product)
}

func (r *DynamoProductRepository) Delete(ctx context.Context, id string) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#Delete")
	defer span.End()
	span.SetAttributes(
		tracing.StringAttribute("productId", id),
	)

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
