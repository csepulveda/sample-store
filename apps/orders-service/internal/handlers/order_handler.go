package handlers

import (
	"orders-service/internal/domain"
	"orders-service/internal/repository"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateOrderHandler handles POST /api/orders
func CreateOrderHandler(repo repository.OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Items []domain.OrderItem `json:"items"`
		}

		if err := c.BodyParser(&input); err != nil {
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
		}

		if err := repo.Create(&order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(order)
	}
}

// ListOrdersHandler handles GET /api/orders
func ListOrdersHandler(repo repository.OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if id != "" {
			order, err := repo.GetByID(id)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Order not found",
				})
			}
			return c.JSON(order)
		}

		orders, err := repo.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(orders)
	}
}

// PatchOrderHandler handles PATCH /api/orders/:id
func PatchOrderHandler(repo repository.OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		order, err := repo.GetByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}

		patchData := make(map[string]interface{})
		if err := c.BodyParser(&patchData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if status, ok := patchData["status"].(string); ok {
			order.Status = status
		}

		if err := repo.Update(order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(order)
	}
}

// DeleteOrderHandler handles DELETE /api/orders/:id
func DeleteOrderHandler(repo repository.OrderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := repo.Delete(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
