package payment_controller

import (
	"errors"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	payment_usecases "github.com/Rabi-IT/rabi-food-core/features/payment/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type chargeBody struct {
	SubscriptionID    string `json:"subscriptionId"    validate:"required"`
	PaymentToken string `json:"paymentToken"`
}

type chargeResponse struct {
	RequiresAction bool   `json:"requiresAction,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
}

func (c *PaymentController) Charge(ctx *fiber.Ctx) error {
	var body chargeBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("invalid body")
	}

	if err := validator.V.Struct(body); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)
	logger.GetWideEvent(uctx).Event = "payment-charge"

	out, err := c.usecase.Charge(uctx, payment_usecases.ChargeInput{
		SubscriptionID:    body.SubscriptionID,
		PaymentToken: body.PaymentToken,
		UserID:            session.UserID,
		UserEmail:         session.Login,
		UserName:          session.Name,
		TenantID:          ctx.Get("X-Tenant-ID"),
	})

	if errors.Is(err, errs.ErrPaymentRequiresAction) {
		return ctx.Status(http.StatusOK).JSON(chargeResponse{
			RequiresAction: true,
			ClientSecret:   out.ClientSecret,
		})
	}

	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
