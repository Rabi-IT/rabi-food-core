package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-delete"

	if err := c.usecase.Delete(uctx, id); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
