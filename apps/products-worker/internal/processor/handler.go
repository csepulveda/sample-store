package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"products-worker/internal/repository"
	"products-worker/internal/tracing"
)

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderMessage struct {
	Type     string      `json:"type"`
	OrderID  string      `json:"orderId"`
	Items    []OrderItem `json:"items"`
	Datetime string      `json:"datetime"`
}

type Handler interface {
	HandleMessage(ctx context.Context, message string) error
}

type OrderHandler struct {
	repo repository.ProductRepository
}

func NewOrderHandler(repo repository.ProductRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) HandleMessage(ctx context.Context, message string) error {
	ctx, span := tracing.NewSpan(ctx, "OrderHandler#HandleMessage")
	defer span.End()

	var order OrderMessage
	if err := json.Unmarshal([]byte(message), &order); err != nil {
		return fmt.Errorf("failed to parse order message: %w", err)
	}

	span.SetAttributes(
		tracing.StringAttribute("orderId", order.OrderID),
		tracing.StringAttribute("eventType", order.Type),
	)

	for _, item := range order.Items {
		switch order.Type {
		case "order.created":
			if err := h.repo.DecrementStock(ctx, item.ProductID, item.Quantity); err != nil {
				return fmt.Errorf("failed to decrement stock for product %s: %w", item.ProductID, err)
			}
		case "order.canceled", "order.returned":
			if err := h.repo.IncrementStock(ctx, item.ProductID, item.Quantity); err != nil {
				return fmt.Errorf("failed to increment stock for product %s: %w", item.ProductID, err)
			}
		default:
			return fmt.Errorf("unsupported order event type: %s", order.Type)
		}
	}

	return nil
}
