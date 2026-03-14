package tenant_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GetByID godoc
// @Summary Get tenant by ID
// @Description Get a tenant by ID
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tenant/{id} [get].
func (c *TenantController) GetByID(ctx *fiber.Ctx) error {
	data, err := c.usecase.GetByID(ctx.UserContext(), ctx.Params("id"))

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
