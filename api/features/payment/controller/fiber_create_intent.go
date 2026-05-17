package payment_controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	payment_usecases "github.com/Rabi-IT/rabi-food-core/features/payment/usecases"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type createIntentBody struct {
	Items       []subscription_usecases.SubscriptionItemInput `json:"items"       validate:"required,min=1"`
	TotalCycles uint                                          `json:"totalCycles" validate:"required,min=1"`
}

func (c *PaymentController) CreateIntent(ctx *fiber.Ctx) error {
	var body createIntentBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.JSON(err)
	}

	if err := validator.V.Struct(body); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)
	logger.GetWideEvent(uctx).Event = "payment-create-intent"

	out, err := c.usecase.CreateIntent(uctx, payment_usecases.CreateIntentInput{
		TenantID:    ctx.Get("X-Tenant-ID"),
		UserID:      session.UserID,
		UserEmail:   session.Login,
		UserName:    session.Name,
		Items:       body.Items,
		TotalCycles: body.TotalCycles,
	})
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(out)
}
