package tenant_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/tenant_case"

	"github.com/gofiber/fiber/v2"
)

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
