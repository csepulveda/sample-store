package domain

type Order struct {
	ID        string      `json:"id" dynamodbav:"id"`
	Status    string      `json:"status" dynamodbav:"status"`
	CreatedAt string      `json:"createdAt" dynamodbav:"createdAt"`
	Items     []OrderItem `json:"items" dynamodbav:"items"`
}

type OrderItem struct {
	ProductID string `json:"productId" dynamodbav:"productId"`
	Quantity  int    `json:"quantity" dynamodbav:"quantity"`
}
