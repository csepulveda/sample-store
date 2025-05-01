package repository

import (
	"context"
	"errors"
	"orders-service/internal/domain"
	"orders-service/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoOrderRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoOrderRepository(client *dynamodb.Client, tableName string) *DynamoOrderRepository {
	return &DynamoOrderRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoOrderRepository#Create")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("orderId", order.ID),
	)

	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoOrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	ctx, span := tracing.NewSpan(ctx, "DynamoOrderRepository#GetAll")
	defer span.End()

	output, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	err = attributevalue.UnmarshalListOfMaps(output.Items, &orders)
	if err != nil {
		return nil, err
	}

	// Filtrar eliminadas
	var activeOrders []domain.Order
	for _, o := range orders {
		if !o.Deleted {
			activeOrders = append(activeOrders, o)
		}
	}

	return activeOrders, nil
}

func (r *DynamoOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	ctx, span := tracing.NewSpan(ctx, "DynamoOrderRepository#GetByID")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("orderId", id),
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

	if output.Item == nil {
		return nil, errors.New("order not found")
	}

	var order domain.Order
	err = attributevalue.UnmarshalMap(output.Item, &order)
	return &order, err
}

func (r *DynamoOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoOrderRepository#Update")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("orderId", order.ID),
	)

	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoOrderRepository) Delete(ctx context.Context, id string) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoOrderRepository#Update")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("orderId", id),
	)

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
