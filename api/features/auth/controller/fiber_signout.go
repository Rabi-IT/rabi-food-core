package controller

import (
	"net/http"
	"strings"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) SignOut(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "auth-signout"

	authHeader := ctx.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	if err := c.usecase.SignOut(uctx, accessToken); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
