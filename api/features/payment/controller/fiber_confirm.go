package payment_controller

import (
	"errors"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	payment_usecases "github.com/Rabi-IT/rabi-food-core/features/payment/usecases"
	subscription_gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type confirmBody struct {
	PaymentIntentID string                                        `json:"paymentIntentId"`
	Items           []subscription_usecases.SubscriptionItemInput `json:"items"           validate:"required,min=1"`
	DeliveryDays    []subscription_gateway.DeliveryDay            `json:"deliveryDays"    validate:"required,min=1"`
	TotalCycles     uint                                          `json:"totalCycles"     validate:"required,min=1"`
	Notes           string                                        `json:"notes"`
}

type confirmResponse struct {
	SubscriptionIDs []string `json:"subscriptionIds,omitempty"`
	ClientSecret    string   `json:"clientSecret,omitempty"`
	RequiresAction  bool     `json:"requiresAction,omitempty"`
}

func (c *PaymentController) Confirm(ctx *fiber.Ctx) error {
	var body confirmBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.JSON(err)
	}

	if err := validator.V.Struct(body); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)
	logger.GetWideEvent(uctx).Event = "payment-confirm"

	out, err := c.usecase.Confirm(uctx, payment_usecases.ConfirmInput{
		TenantID:        ctx.Get("X-Tenant-ID"),
		UserID:          session.UserID,
		UserEmail:       session.Login,
		UserName:        session.Name,
		PaymentIntentID: body.PaymentIntentID,
		Items:           body.Items,
		DeliveryDays:    body.DeliveryDays,
		TotalCycles:     body.TotalCycles,
		Notes:           body.Notes,
	})

	if errors.Is(err, errs.ErrPaymentRequiresAction) {
		return ctx.Status(http.StatusOK).JSON(confirmResponse{
			ClientSecret:   out.ClientSecret,
			RequiresAction: true,
		})
	}

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(confirmResponse{SubscriptionIDs: out.SubscriptionIDs})
}
