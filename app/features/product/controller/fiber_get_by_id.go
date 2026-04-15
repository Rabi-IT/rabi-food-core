package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"

	"github.com/gofiber/fiber/v2"
)

// keep import for swagger documentation.
var _ *gateway.GetByIDOutput

// GetByID godoc
// @Summary Get product by ID
// @Description Get a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/{id} [get].
func (c *ProductController) GetByID(ctx *fiber.Ctx) error {
	filter := gateway.GetByIDFilter{
		ID:       ctx.Params("id"),
		TenantID: ctx.Get("X-Tenant-ID"),
	}

	data, err := c.usecase.GetByID(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
