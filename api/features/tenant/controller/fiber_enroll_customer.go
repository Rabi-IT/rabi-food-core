package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) EnrollCustomer(ctx *fiber.Ctx) error {
	err := c.usecase.EnrollCustomer(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
