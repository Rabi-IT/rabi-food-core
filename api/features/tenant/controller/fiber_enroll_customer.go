package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) EnrollCustomer(ctx *fiber.Ctx) error {
	if err := c.usecase.EnrollCustomer(ctx.UserContext(), ctx.Params("id")); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
