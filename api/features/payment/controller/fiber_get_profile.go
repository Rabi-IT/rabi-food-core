package payment_controller

import (
	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/gofiber/fiber/v2"
)

func (c *PaymentController) GetProfile(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "payment-get-profile"

	userID := app_context.GetSession(uctx).UserID

	out, err := c.usecase.GetProfile(uctx, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(out)
}
