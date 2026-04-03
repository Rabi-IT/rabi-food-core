package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/http/parser"

	"github.com/gofiber/fiber/v2"
)

// Patch godoc
// @Summary Patch tenant
// @Description Update an existing tenant
// @Tags tenants
// @Accept json
// @Produce text/plain
// @Param id path string true "Tenant ID"
// @Param tenant body usecases.PatchValues true "Tenant patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tenant/{id} [patch].
func (c *TenantController) Patch(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	filter := usecases.PatchFilter{
		ID: &id,
	}

	data := usecases.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
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
