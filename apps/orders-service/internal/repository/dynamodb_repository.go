package repository

import (
	"context"
	"errors"
	"orders-service/internal/domain"

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

func (r *DynamoOrderRepository) Create(order *domain.Order) error {
	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoOrderRepository) GetAll() ([]domain.Order, error) {
	output, err := r.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	err = attributevalue.UnmarshalListOfMaps(output.Items, &orders)
	return orders, err
}

func (r *DynamoOrderRepository) GetByID(id string) (*domain.Order, error) {
	output, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

func (r *DynamoOrderRepository) Update(order *domain.Order) error {
	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoOrderRepository) Delete(id string) error {
	_, err := r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
