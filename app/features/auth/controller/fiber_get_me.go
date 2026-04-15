package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) GetMe(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-get-me"

	id := app_context.GetSession(uctx).UserID
	out, err := c.usecase.GetByID(uctx, id)
	if err != nil {
		return err
	}
	if out == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(out)
}
