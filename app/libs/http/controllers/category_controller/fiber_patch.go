package category_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// PatchCategory godoc
// @Summary Patch category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce text/plain
// @Param id path string true "Category ID"
// @Param category body category_gateway.PatchValues true "Category patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /category/{id} [patch]
func (c *CategoryController) Patch(ctx *fiber.Ctx) error {
	filter := category_gateway.PatchFilter{
		ID: ctx.Params("id"),
	}

	data := category_gateway.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	updated, err := c.usecase.Patch(ctx.UserContext(), filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
