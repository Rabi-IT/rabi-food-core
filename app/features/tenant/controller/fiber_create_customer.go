package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) CreateCustomer(ctx *fiber.Ctx) error {
	if err := c.usecase.CreateCustomer(ctx.UserContext(), ctx.Params("id")); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
