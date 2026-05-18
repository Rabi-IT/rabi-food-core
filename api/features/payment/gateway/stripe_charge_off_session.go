package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

func (g *StripePaymentGatewayAdapter) ChargeOffSession(ctx context.Context, input ChargeOffSessionInput) (ChargeOutput, error) {
	user, err := g.auth.GetByID(ctx, input.UserID)
	if err != nil {
		return ChargeOutput{}, fmt.Errorf("failed to look up customer: %w", err)
	}

	if user == nil || user.StripeCustomerID == "" || user.StripePaymentMethodID == "" {
		return ChargeOutput{}, errs.ErrNoPaymentMethod
	}

	params := &stripe.PaymentIntentCreateParams{
		Amount:        stripe.Int64(input.AmountCents),
		Currency:      stripe.String(input.Currency),
		Customer:      stripe.String(user.StripeCustomerID),
		PaymentMethod: stripe.String(user.StripePaymentMethodID),
		Confirm:       stripe.Bool(true),
		OffSession:    stripe.Bool(true),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Metadata: input.Metadata,
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
