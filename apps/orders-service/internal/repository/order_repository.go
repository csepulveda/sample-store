package repository

import "orders-service/internal/domain"

type OrderRepository interface {
	Create(order *domain.Order) error
	GetAll() ([]domain.Order, error)
	GetByID(id string) (*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id string) error
}
