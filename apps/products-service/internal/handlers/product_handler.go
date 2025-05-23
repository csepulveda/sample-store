package handlers

import (
	"products-service/internal/domain"
	"products-service/internal/repository"
	"products-service/internal/tracing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateProductHandler handles POST /api/products
func CreateProductHandler(repo repository.ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "CreateProductHandler")
		defer span.End()

		var product domain.Product

		if err := c.BodyParser(&product); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Generate UUID for new product
		product.ID = uuid.New().String()

		if err := repo.Create(ctx, &product); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}

// ListProductsHandler handles GET /api/products
func ListProductsHandler(repo repository.ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "ListProductsHandler")
		defer span.End()
		id := c.Params("id")
		span.SetAttributes(
			tracing.StringAttribute("productId", id),
		)

		if id != "" {
			product, err := repo.GetByID(ctx, id)
			if err != nil {
				span.RecordError(err)
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Product not found",
				})
			}
			return c.JSON(product)
		}

		products, err := repo.GetAll(ctx)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(products)
	}
}

// UpdateProductHandler handles PUT /api/products/:id
func UpdateProductHandler(repo repository.ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "UpdateProductHandler")
		defer span.End()

		id := c.Params("id")
		span.SetAttributes(
			tracing.StringAttribute("productId", id),
		)

		product, err := repo.GetByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}

		if err := c.BodyParser(product); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		product.ID = id // aseguramos que no se modifique el ID

		if err := repo.Update(ctx, product); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(product)
	}
}

// PatchProductHandler handles PATCH /api/products/:id
func PatchProductHandler(repo repository.ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "PatchProductHandler")
		defer span.End()
		id := c.Params("id")

		span.SetAttributes(
			tracing.StringAttribute("productId", id),
		)

		product, err := repo.GetByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}

		patchData := make(map[string]interface{})
		if err := c.BodyParser(&patchData); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if name, ok := patchData["name"].(string); ok {
			product.Name = name
		}
		if description, ok := patchData["description"].(string); ok {
			product.Description = description
		}
		if price, ok := patchData["price"].(float64); ok {
			product.Price = price
		}
		if stock, ok := patchData["stock"].(float64); ok {
			product.Stock = int(stock)
		}

		if err := repo.Update(ctx, product); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(product)
	}
}

// DeleteProductHandler handles DELETE /api/products/:id
func DeleteProductHandler(repo repository.ProductRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracing.NewSpan(c.UserContext(), "DeleteProductHandler")
		defer span.End()
		id := c.Params("id")

		span.SetAttributes(
			tracing.StringAttribute("productId", id),
		)

		if err := repo.Delete(ctx, id); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
