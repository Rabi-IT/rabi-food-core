package product_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *ProductController) Patch(ctx *fiber.Ctx) error {
	filter := product_gateway.PatchFilter{
		ID: ctx.Params("id"),
	}

	data := product_gateway.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
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
