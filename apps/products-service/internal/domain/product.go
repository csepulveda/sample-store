package domain

type Product struct {
	ID          string  `json:"id" dynamodbav:"id"`
	Name        string  `json:"name" dynamodbav:"name"`
	Description string  `json:"description" dynamodbav:"description"`
	Price       float64 `json:"price" dynamodbav:"price"`
	Stock       int     `json:"stock" dynamodbav:"stock"`
}
