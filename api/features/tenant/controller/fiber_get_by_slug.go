package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) GetBySlug(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "tenant-get-by-slug"

	out, err := c.usecase.GetBySlug(uctx, ctx.Params("slug"))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(out)
}
