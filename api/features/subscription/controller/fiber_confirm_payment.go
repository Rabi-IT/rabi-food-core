package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (c *SubscriptionController) OnPaymentReceived(ctx *fiber.Ctx) error {
	sig := ctx.Get("Stripe-Signature")
	body := ctx.Body()

	event, err := webhook.ConstructEvent(body, sig, config.StripeWebhookSecret)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("invalid signature")
	}

	if event.Type != stripe.EventTypePaymentIntentSucceeded {
		return ctx.SendStatus(http.StatusOK)
	}

	subscriptionID := event.GetObjectValue("metadata", "subscription_id")
	if subscriptionID == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	amountStr := event.GetObjectValue("amount")
	amountCents, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("invalid amount")
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "subscription-payment-received"

	err = c.usecase.OnPaymentReceived(uctx, usecases.OnPaymentReceivedInput{
		SubscriptionID:  subscriptionID,
		PaymentIntentID: event.GetObjectValue("id"),
		AmountCents:     amountCents,
		Provider:        "stripe",
		PaidAt:          time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
