package usecases

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type OnPaymentReceivedInput struct {
	SubscriptionID  string
	PaymentIntentID string
	AmountCents     int64
	Provider        string
	PaidAt          time.Time
}

func (c *SubscriptionCase) OnPaymentReceived(ctx context.Context, in OnPaymentReceivedInput) error {
	wd := logger.GetWideEvent(ctx)
	wd.Event = "subscription-on-payment-received"

	sub, err := c.gateway.GetByID(ctx, g.GetByIDFilter{ID: in.SubscriptionID})
	if err != nil {
		return err
	}

	wd.TrackSubscriptionPayment(in.SubscriptionID, in.PaymentIntentID, in.Provider, in.AmountCents)

	if sub == nil || int64(sub.PaymentAmount) != in.AmountCents {
		return c.paymentGateway.Refund(ctx, in.PaymentIntentID)
	}

	paidStatus := payment_status.StatusPaid

	updated, err := c.gateway.Patch(ctx, g.PatchFilter{
		ID:            in.SubscriptionID,
		PaymentStatus: payment_status.StatusPending,
	}, g.PatchValues{
		PaymentStatus: &paidStatus,
		PaidAt:        &in.PaidAt,
	})
	if err != nil {
		return err
	}

	if updated {
		return nil
	}

	found, err := c.gateway.GetByID(ctx, g.GetByIDFilter{ID: in.SubscriptionID})
	if err != nil {
		return err
	}

	if found != nil && found.PaymentStatus == payment_status.StatusPaid {
		wd.MarkPaymentIdempotent()
		return nil
	}

	return errs.ErrSubscriptionStateVerificationFailed
}
