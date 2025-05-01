package handlers

import (
	"orders-service/internal/domain"
	"orders-service/internal/publisher"
	"orders-service/internal/repository"
	"orders-service/internal/tracing"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateOrderHandler handles POST /api/orders
func CreateOrderHandler(repo repository.OrderRepository, pub publisher.OrderPublisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "CreateOrderHandler")
		defer span.End()

		var input struct {
			Items []domain.OrderItem `json:"items"`
		}

		if err := c.BodyParser(&input); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if len(input.Items) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Order must have at least one item",
			})
		}

		order := domain.Order{
			ID:        uuid.New().String(),
			Status:    "created",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			Items:     input.Items,
			Deleted:   false,
		}

		if err := repo.Create(ctx, &order); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := pub.PublishOrderCreated(ctx, order); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Order created but failed to publish event: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(order)
	}
}

// ListOrdersHandler handles GET /api/orders
func ListOrdersHandler(repo repository.OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "ListOrdersHandler")
		defer span.End()

		id := c.Params("id")

		span.SetAttributes(
			tracing.StringAttribute("orderId", id),
		)

		if id != "" {
			order, err := repo.GetByID(ctx, id)
			if err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Order not found",
				})
			}
			return c.JSON(order)
		}

		orders, err := repo.GetAll(ctx)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(orders)
	}
}

// PatchOrderHandler handles PATCH /api/orders/:id
func PatchOrderHandler(repo repository.OrderRepository, pub publisher.OrderPublisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "PatchOrderHandler")
		defer span.End()

		id := c.Params("id")
		span.SetAttributes(
			tracing.StringAttribute("orderId", id),
		)

		order, err := repo.GetByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}

		if order.Deleted {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot modify a deleted order",
			})
		}

		patchData := make(map[string]interface{})
		if err := c.BodyParser(&patchData); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if status, ok := patchData["status"].(string); ok {
			validTransition := true

			switch order.Status {
			case "canceled", "returned":
				validTransition = false
			case "delivered":
				validTransition = (status == "returned")
			default:
				validTransition = true
			}

			if !validTransition {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid state transition",
				})
			}

			order.Status = status
		}

		if err := repo.Update(ctx, order); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Publish events if necessary
		switch order.Status {
		case "canceled":
			if err := pub.PublishOrderCanceled(ctx, *order); err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Order updated but failed to publish cancel event: " + err.Error(),
				})
			}
		case "returned":
			if err := pub.PublishOrderCanceled(ctx, *order); err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Order updated but failed to publish return event: " + err.Error(),
				})
			}
		}

		return c.JSON(order)
	}
}

// DeleteOrderHandler handles DELETE /api/orders/:id
func DeleteOrderHandler(repo repository.OrderRepository, pub publisher.OrderPublisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "DeleteOrderHandler")
		defer span.End()

		id := c.Params("id")
		span.SetAttributes(
			tracing.StringAttribute("orderId", id),
		)

		order, err := repo.GetByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}

		if order.Status != "delivered" && order.Status != "canceled" && order.Status != "returned" {
			order.Status = "canceled"

			if err := repo.Update(ctx, order); err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to cancel before delete: " + err.Error(),
				})
			}

			if err := pub.PublishOrderCanceled(ctx, *order); err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to publish cancel event: " + err.Error(),
				})
			}
		}

		order.Deleted = true
		if err := repo.Update(ctx, order); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to soft-delete order: " + err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
