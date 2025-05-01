package repository

import (
	"context"
	"fmt"
	"products-worker/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductRepository interface {
	DecrementStock(ctx context.Context, id string, qty int) error
	IncrementStock(ctx context.Context, id string, qty int) error
}

type DynamoProductRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoProductRepository(client *dynamodb.Client, table string) *DynamoProductRepository {
	return &DynamoProductRepository{client: client, tableName: table}
}

func (r *DynamoProductRepository) updateStock(ctx context.Context, id string, qty int) error {
	ctx, span := tracing.NewSpan(ctx, "DynamoProductRepository#updateStock")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("productId", id),
		tracing.IntAttribute("quantity", qty),
	)

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}

	update := "SET stock = stock + :val"
	values := map[string]types.AttributeValue{
		":val": &types.AttributeValueMemberN{Value: stringInt(qty)},
	}

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       key,
		UpdateExpression:          &update,
		ExpressionAttributeValues: values,
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	return err
}

func (r *DynamoProductRepository) DecrementStock(ctx context.Context, id string, qty int) error {
	return r.updateStock(ctx, id, -qty)
}

func (r *DynamoProductRepository) IncrementStock(ctx context.Context, id string, qty int) error {
	return r.updateStock(ctx, id, qty)
}

func stringInt(i int) string {
	return fmt.Sprintf("%d", i)
}
