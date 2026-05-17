package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type ConfirmPaymentInput struct {
	SubscriptionID string
	Provider       string
	PaidAt         time.Time
}

func (c *SubscriptionCase) ConfirmPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	session := app_context.GetSession(ctx)
	if !session.Role.IsSystem() {
		return false, errs.ErrForbidden
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
		if errs.IsUniqueViolation(err) {
			return false, errs.ErrConflict
		}

		return false, err
	}

	if updated {
		return true, nil
	}

	return c.handleNotConfirmedSubscriptionPayment(ctx, in)
}

func (c *SubscriptionCase) handleNotConfirmedSubscriptionPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	wd := logger.GetWideEvent(ctx)

	found, err := c.gateway.GetByID(ctx, g.GetByIDFilter{
		ID: in.SubscriptionID,
	})
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrSubscriptionStateVerificationFailed, err)
	}

	if found == nil {
		return false, errs.ErrSubscriptionNotFound
	}

	if found.PaymentStatus == payment_status.StatusPaid {
		wd.MarkPaymentIdempotent()

		return true, nil
	}

	if found.PaymentStatus != payment_status.StatusPending {
		return false, errs.InvalidTranstion(found.PaymentStatus, payment_status.StatusPaid)
	}

	return false, errs.ErrUnknownBusinessReason
}
