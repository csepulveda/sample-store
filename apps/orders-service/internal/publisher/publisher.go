package publisher

import (
	"context"
	"orders-service/internal/domain"
)

type OrderPublisher interface {
	PublishOrderCreated(ctx context.Context, order domain.Order) error
	PublishOrderCanceled(ctx context.Context, order domain.Order) error
}
