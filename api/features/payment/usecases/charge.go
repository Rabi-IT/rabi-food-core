package payment_usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
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

	metadata := map[string]string{
		"subscription_id": input.SubscriptionID,
		"tenant_id":       input.TenantID,
	}
	idempotencyKey := "charge-" + input.SubscriptionID

	var result g.ChargeOutput

	if input.PaymentToken != "" {
		result, err = c.payment.ChargeWithToken(ctx, g.ChargeWithTokenInput{
			UserID:         input.UserID,
			UserEmail:      input.UserEmail,
			UserName:       input.UserName,
			PaymentToken:   input.PaymentToken,
			AmountCents:    int64(sub.PaymentAmount),
			Currency:       paymentCurrency,
			Metadata:       metadata,
			IdempotencyKey: idempotencyKey,
		})
	} else {
		result, err = c.payment.ChargeOffSession(ctx, g.ChargeOffSessionInput{
			UserID:         input.UserID,
			AmountCents:    int64(sub.PaymentAmount),
			Currency:       paymentCurrency,
			Metadata:       metadata,
			IdempotencyKey: idempotencyKey,
		})
	}

	if err != nil {
		return ChargeOutput{}, err
	}

	if result.Status == "requires_action" {
		return ChargeOutput{RequiresAction: true, ClientSecret: result.ClientSecret}, errs.ErrPaymentRequiresAction
	}

	return ChargeOutput{}, nil
}
