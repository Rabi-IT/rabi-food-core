package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) Patch(ctx *fiber.Ctx) error {
	data := &usecases.PatchInput{}
	if err := ctx.BodyParser(data); err != nil {
		return ctx.JSON(err)
	}

	if err := validator.V.Struct(data); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-patch"

	id := app_context.GetSession(uctx).UserID
	ok, err := c.usecase.Patch(uctx, id, data)
	if err != nil {
		return err
	}
	if !ok {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.SendStatus(http.StatusNoContent)
}
