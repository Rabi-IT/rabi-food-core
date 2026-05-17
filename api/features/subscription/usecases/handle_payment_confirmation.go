package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

// PaymentRefunder is satisfied by the Stripe gateway adapter.
type PaymentRefunder interface {
	Refund(ctx context.Context, paymentIntentID string) error
}

type HandlePaymentConfirmationInput struct {
	SubscriptionID  string
	PaymentIntentID string
	AmountCents     int64
	Provider        string
	PaidAt          time.Time
}

func (c *SubscriptionCase) HandlePaymentConfirmation(ctx context.Context, in HandlePaymentConfirmationInput) error {
	wd := logger.GetWideEvent(ctx)
	wd.Event = "subscription-handle-payment-confirmation"

	sub, err := c.gateway.GetByID(ctx, g.GetByIDFilter{ID: in.SubscriptionID})
	if err != nil {
		return err
	}

	if sub == nil {
		wd.TrackSubscriptionPayment(in.SubscriptionID, in.PaymentIntentID, in.Provider, in.AmountCents)
		return c.paymentRefunder.Refund(ctx, in.PaymentIntentID)
	}

	if int64(sub.PaymentAmount) != in.AmountCents {
		wd.TrackSubscriptionPayment(in.SubscriptionID, in.PaymentIntentID, in.Provider, in.AmountCents)
		return c.paymentRefunder.Refund(ctx, in.PaymentIntentID)
	}

	wd.TrackSubscriptionPayment(in.SubscriptionID, in.PaymentIntentID, in.Provider, in.AmountCents)

	sysCtx := app_context.WithSession(ctx, &app_context.UserSession{Role: auth.System})

	_, err = c.ConfirmPayment(sysCtx, ConfirmPaymentInput{
		SubscriptionID: in.SubscriptionID,
		Provider:       in.Provider,
		PaidAt:         in.PaidAt,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrSubscriptionStateVerificationFailed, err)
	}

	return nil
}
