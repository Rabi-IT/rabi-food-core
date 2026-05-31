package payment_usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

const paymentCurrency = "brl"

type ChargeInput struct {
	SubscriptionID string
	UserID         string
	UserEmail      string
	UserName       string
	TenantID       string
	PaymentToken   string // empty → charge saved card
}

type ChargeOutput struct {
	RequiresAction bool
	ClientSecret   string
}

func (c *PaymentCase) Charge(ctx context.Context, input ChargeInput) (ChargeOutput, error) {
	sub, err := c.subscriptionCase.GetByID(ctx, subscription_gateway.GetByIDFilter{ID: input.SubscriptionID})
	if err != nil {
		return ChargeOutput{}, err
	}

	if sub == nil {
		return ChargeOutput{}, errs.ErrSubscriptionNotFound
	}

	if sub.PaymentStatus != payment_status.StatusPending {
		return ChargeOutput{}, errs.ErrPaymentAlreadyProcessed
	}

	idempotencyKey := "charge-" + input.SubscriptionID

	var (
		result     g.ChargeOutput
		chargeType string
	)

	if input.PaymentToken != "" {
		chargeType = "token"
		result, err = c.payment.ChargeWithToken(ctx, g.ChargeWithTokenInput{
			UserID:         input.UserID,
			UserEmail:      input.UserEmail,
			UserName:       input.UserName,
			PaymentToken:   input.PaymentToken,
			SubscriptionID: input.SubscriptionID,
			TenantID:       input.TenantID,
			AmountCents:    sub.PaymentAmount,
			Currency:       paymentCurrency,
			IdempotencyKey: idempotencyKey,
		})
	} else {
		chargeType = "off_session"
		result, err = c.payment.ChargeOffSession(ctx, g.ChargeOffSessionInput{
			UserID:         input.UserID,
			SubscriptionID: input.SubscriptionID,
			TenantID:       input.TenantID,
			AmountCents:    sub.PaymentAmount,
			Currency:       paymentCurrency,
			IdempotencyKey: idempotencyKey,
		})
	}

	if err != nil {
		return ChargeOutput{}, err
	}

	requiresAction := result.Status == "requires_action"
	logger.GetWideEvent(ctx).TrackCharge(chargeType, result.Status, requiresAction)

	if requiresAction {
		return ChargeOutput{RequiresAction: true, ClientSecret: result.ClientSecret}, errs.ErrPaymentRequiresAction
	}

	return ChargeOutput{}, nil
}
