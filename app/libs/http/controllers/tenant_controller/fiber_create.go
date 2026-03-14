package tenant_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/tenant_case"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create tenant
// @Description Create a new tenant and its first user
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body tenant_case.CreateInput true "Tenant data"
// @Success 201 {object} tenant_case.CreateOutput
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /tenant [post].
func (c *TenantController) Create(ctx *fiber.Ctx) error {
	data := tenant_case.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	output, err := c.usecase.Create(ctx.UserContext(), data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(output)
}
