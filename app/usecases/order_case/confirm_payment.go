package order_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/payment"
	g "rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"
	"time"
)

type ConfirmPaymentInput struct {
	OrderID           string
	ExternalPaymentID string    `json:"externalPaymentId" validate:"required"`
	Provider          string    `json:"provider"          validate:"required"`
	PaidAt            time.Time `json:"paidAt"            validate:"required"`
}

func (c *OrderCase) ConfirmPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	l := logger.Get(ctx).With().
		Str(logger.OrderID, in.OrderID).
		Str(logger.PaymentID, in.ExternalPaymentID).
		Logger()

	session := app_context.GetSession(ctx)
	if !session.Role.IsSystem() {
		return false, errs.ErrForbidden
	}

	updated, err := c.gateway.Patch(g.PatchFilter{
		ID:            in.OrderID,
		PaymentStatus: payment.StatusPending,
	}, g.PatchValues{
		ExternalPaymentID: in.ExternalPaymentID,
		PaymentStatus:     payment.StatusPaid,
		Provider:          in.Provider,
		PaidAt:            in.PaidAt,
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

	l.Warn().Msg("payment confirmation not updated, checking order status")

	ctx = logger.WithContext(ctx, l)

	return c.handleNotConfirmedPayment(ctx, in)
}

func (c *OrderCase) handleNotConfirmedPayment(ctx context.Context, in ConfirmPaymentInput) (bool, error) {
	l := logger.Get(ctx)

	orderFound, err := c.gateway.GetByID(g.GetByIDFilter{ID: in.OrderID})
	if err != nil {
		return false, err
	}

	if orderFound == nil {
		l.Error().Msg("order not found when confirming payment")

		return false, nil
	}

	if orderFound.PaymentStatus == payment.StatusPaid {
		if orderFound.ExternalPaymentID != nil && *orderFound.ExternalPaymentID == in.ExternalPaymentID {
			l.Info().Msg("payment already confirmed for this order")

			return true, nil
		}

		l.Error().Msg("conflict on confirming payment, external_payment_id already used")

		return false, errs.ErrConflict
	}

	if orderFound.PaymentStatus != payment.StatusPending {
		l.Error().
			Str("current_status", string(orderFound.PaymentStatus)).
			Msg("invalid transition when confirming payment")

		return false, errs.InvalidTranstion(orderFound.PaymentStatus, payment.StatusPaid)
	}

	l.Error().Msg("payment confirmation not updated due to unknown reason")

	return false, nil
}
