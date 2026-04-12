package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) RegisterCustomer(ctx *fiber.Ctx) error {
	if err := c.usecase.RegisterCustomer(ctx.UserContext(), ctx.Params("id")); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
