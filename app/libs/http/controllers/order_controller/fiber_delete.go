package order_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"

	"github.com/gofiber/fiber/v2"
)

func (c *OrderController) Delete(ctx *fiber.Ctx) error {
	filter := order_gateway.DeleteFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return ctx.JSON(err)
	}

	filter.ID = ctx.Params("id")
	deleted, err := c.usecase.Delete(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
