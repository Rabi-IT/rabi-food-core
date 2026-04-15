package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	parser "github.com/Rabi-IT/rabi-food-core/libs/http/parser"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Patch godoc
// @Summary Patch category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce text/plain
// @Param id path string true "Category ID"
// @Param category body gateway.PatchValues true "Category patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /category/{id} [patch].
func (c *CategoryController) Patch(ctx *fiber.Ctx) error {
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

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "patch-category"
	updated, err := c.usecase.Patch(uctx, filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
