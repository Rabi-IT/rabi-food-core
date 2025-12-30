package category_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	data := category_gateway.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	id, err := c.usecase.Create(ctx.UserContext(), data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
