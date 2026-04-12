package controller

import (
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-get-by-id"

	out, err := c.usecase.GetByID(uctx, id)
	if err != nil {
		return err
	}

	return ctx.JSON(out)
}

func (c *AuthController) GetMe(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-get-me"

	out, err := c.usecase.GetMe(uctx)
	if err != nil {
		return err
	}

	return ctx.JSON(out)
}
