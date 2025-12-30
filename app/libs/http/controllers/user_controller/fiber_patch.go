package user_controller

import (
	"net/http"
	g "rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) Patch(ctx *fiber.Ctx) error {
	data := g.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	updated, err := c.usecase.Patch(ctx.UserContext(), ctx.Params("id"), data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
