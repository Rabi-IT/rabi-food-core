package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create tenant
// @Description Create a new tenant and its first user
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body usecases.CreateInput true "Tenant data"
// @Success 201 {object} usecases.CreateOutput
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /tenant [post].
func (c *TenantController) Create(ctx *fiber.Ctx) error {
	data := usecases.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "new-tenant"
	output, err := c.usecase.Create(uctx, data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(output)
}
