package product_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	data := product_gateway.CreateInput{}
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
