package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"

	"github.com/gofiber/fiber/v2"
)

// keep import for swagger documentation.
var _ *gateway.GetByIDOutput

// GetByID godoc
// @Summary Get category by ID
// @Description Get a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /category/{id} [get].
func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	filter := gateway.GetByIDFilter{
		ID: ctx.Params("id"),
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
