package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"

	"github.com/gofiber/fiber/v2"
)

// keep import for swagger documentation.
var _ *gateway.GetByIDOutput

// GetByID godoc
// @Summary Get tenant by ID
// @Description Get a tenant by ID
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tenant/{id} [get].
func (c *TenantController) GetByID(ctx *fiber.Ctx) error {
	data, err := c.usecase.GetByID(ctx.UserContext(), ctx.Params("id"), ctx.Get("X-Tenant-ID"))

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
