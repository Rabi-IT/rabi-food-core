package order_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *OrderController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.usecase.GetByID(ctx.UserContext(), id)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
