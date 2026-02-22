package category_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"

	"github.com/gofiber/fiber/v2"
)

func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	filter := category_gateway.DeleteFilter{}

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
