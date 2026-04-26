package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) VerifyOTP(ctx *fiber.Ctx) error {
	data := &usecases.VerifyOTPInput{}
	if err := ctx.BodyParser(data); err != nil {
		return ctx.JSON(err)
	}

	if err := validator.V.Struct(data); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "auth-verify-otp"

	out, err := c.usecase.VerifyOTP(uctx, data)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(out)
}
