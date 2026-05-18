package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

const stripeProvider = "stripe"

func (c *SubscriptionController) OnPaymentReceived(ctx *fiber.Ctx) error {
	event, err := webhook.ConstructEvent(ctx.Body(), ctx.Get("Stripe-Signature"), config.StripeWebhookSecret)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidWebhookSignature, err)
	}

	// Stripe requires 200 for unhandled event types, otherwise it retries indefinitely.
	if event.Type != stripe.EventTypePaymentIntentSucceeded {
		return ctx.SendStatus(http.StatusOK)
	}

	amountCents, err := strconv.ParseInt(event.GetObjectValue("amount"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrMalformedWebhookPayload, err)
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "subscription-payment-received"

	err = c.usecase.OnPaymentReceived(uctx, usecases.OnPaymentReceivedInput{
		SubscriptionID:  event.GetObjectValue("metadata", "subscription_id"),
		PaymentIntentID: event.GetObjectValue("id"),
		AmountCents:     amountCents,
		Provider:        stripeProvider,
		PaidAt:          time.Now().UTC(),
	})

	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
