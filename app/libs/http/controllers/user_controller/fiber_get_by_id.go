package user_controller

import (
	"net/http"
	g "rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	filter := g.GetByIDFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return err
	}

	filter.ID = ctx.Params("id")
	data, err := c.usecase.GetByID(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
