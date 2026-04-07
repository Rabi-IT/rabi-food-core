package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type ConfirmPaymentInput struct {
	OrderID           string
	ExternalPaymentID string    `json:"externalPaymentId" validate:"required"`
	Provider          string    `json:"provider"          validate:"required"`
	PaidAt            time.Time `json:"paidAt"            validate:"required"`
}

func (c *OrderCase) ConfirmPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	wd := logger.GetWideEvent(ctx)
	wd.TrackOrderPayment(in.OrderID, in.ExternalPaymentID)

	session := app_context.GetSession(ctx)
	if !session.Role.IsSystem() {
		return false, errs.ErrForbidden
	}

	paidStatus := payment_status.StatusPaid
	updated, err := c.gateway.Patch(g.PatchFilter{
		ID:            in.OrderID,
		PaymentStatus: payment_status.StatusPending,
	}, g.PatchValues{
		ExternalPaymentID: &in.ExternalPaymentID,
		PaymentStatus:     &paidStatus,
		Provider:          &in.Provider,
		PaidAt:            &in.PaidAt,
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

	return c.handleNotConfirmedPayment(ctx, in)
}

func (c *OrderCase) handleNotConfirmedPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	wd := logger.GetWideEvent(ctx)

	orderFound, err := c.gateway.GetByID(g.GetByIDFilter{ID: in.OrderID})
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrOrderStateVerificationFailed, err)
	}

	if orderFound == nil {
		return false, errs.ErrOrderNotFound
	}

	if orderFound.PaymentStatus == payment_status.StatusPaid {
		if orderFound.ExternalPaymentID != nil && *orderFound.ExternalPaymentID == in.ExternalPaymentID {
			wd.MarkPaymentIdempotent()
			return true, nil
		}

		wd.TrackPaymentConflict(*orderFound.ExternalPaymentID)
		return false, errs.ErrPaymentExternalIDConflict
	}

	if orderFound.PaymentStatus != payment_status.StatusPending {
		return false, errs.InvalidTranstion(orderFound.PaymentStatus, payment_status.StatusPaid)
	}

	return false, errs.ErrUnknownBusinessReason
}
