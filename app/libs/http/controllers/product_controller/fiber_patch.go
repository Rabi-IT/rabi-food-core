package product_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// PatchProduct godoc
// @Summary Patch product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce text/plain
// @Param id path string true "Product ID"
// @Param product body product_gateway.PatchValues true "Product patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/{id} [patch]
func (c *ProductController) Patch(ctx *fiber.Ctx) error {
	filter := product_gateway.PatchFilter{
		ID: ctx.Params("id"),
	}

	data := product_gateway.PatchValues{}
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
