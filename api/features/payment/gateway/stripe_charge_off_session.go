package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

var card = "card"

func (g *StripePaymentGatewayAdapter) ChargeOffSession(ctx context.Context, input ChargeOffSessionInput) (ChargeOutput, error) {
	user, err := g.auth.GetByID(ctx, input.UserID)
	if err != nil {
		return ChargeOutput{}, fmt.Errorf("failed to look up customer: %w", err)
	}

	if user == nil || user.StripeCustomerID == "" || user.StripePaymentMethodID == "" {
		return ChargeOutput{}, errs.ErrNoPaymentMethod
	}

	params := &stripe.PaymentIntentCreateParams{
		Amount:        new(int64(input.AmountCents)),
		Currency:      &input.Currency,
		Customer:      &user.StripeCustomerID,
		PaymentMethod: &user.StripePaymentMethodID,
		Confirm:       new(true),
		OffSession:    new(true),
		PaymentMethodTypes: []*string{
			&card,
		},
		Metadata: map[string]string{
			"subscription_id": input.SubscriptionID,
			"tenant_id":       input.TenantID,
		},
	}

	params.SetIdempotencyKey(input.IdempotencyKey)

	pi, err := g.sc.V1PaymentIntents.Create(ctx, params)
	if err != nil {
		return ChargeOutput{}, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	out := ChargeOutput{
		PaymentIntentID: pi.ID,
		Status:          string(pi.Status),
	}

	if pi.Status == stripe.PaymentIntentStatusRequiresAction {
		out.ClientSecret = pi.ClientSecret
	}

	return out, nil
}
