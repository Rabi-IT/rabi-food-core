package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/http/parser"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Patch godoc
// @Summary Patch product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce text/plain
// @Param id path string true "Product ID"
// @Param product body gateway.PatchValues true "Product patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/{id} [patch].
func (c *ProductController) Patch(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()

	filter := gateway.PatchFilter{
		ID:       ctx.Params("id"),
		TenantID: ctx.Get("X-Tenant-ID"),
	}

	data := gateway.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	logger.GetWideEvent(uctx).Event = "patch-product"
	updated, err := c.usecase.Patch(uctx, filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
