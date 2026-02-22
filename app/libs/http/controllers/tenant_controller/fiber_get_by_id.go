package tenant_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) GetByID(ctx *fiber.Ctx) error {
	data, err := c.usecase.GetByID(ctx.UserContext(), ctx.Params("id"))

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
