package product_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	deleted, err := c.usecase.Delete(ctx.UserContext(), id)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
