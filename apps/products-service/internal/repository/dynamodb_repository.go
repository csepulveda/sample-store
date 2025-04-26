package repository

import (
	"context"
	"errors"
	"products-service/internal/domain"

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

func (r *DynamoProductRepository) Create(product *domain.Product) error {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoProductRepository) GetAll() ([]domain.Product, error) {
	output, err := r.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	err = attributevalue.UnmarshalListOfMaps(output.Items, &products)
	return products, err
}

func (r *DynamoProductRepository) GetByID(id string) (*domain.Product, error) {
	output, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

func (r *DynamoProductRepository) Update(product *domain.Product) error {
	return r.Create(product)
}

func (r *DynamoProductRepository) Delete(id string) error {
	_, err := r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
