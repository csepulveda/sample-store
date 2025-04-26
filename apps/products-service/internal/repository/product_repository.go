package repository

import "products-service/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id string) error
}
