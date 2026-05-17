package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

func (g *StripeGatewayAdapter) ConfirmOffSession(ctx context.Context, input ConfirmOffSessionInput) (ConfirmOffSessionOutput, error) {
	createParams := &stripe.PaymentIntentCreateParams{
		Amount:        &input.AmountCents,
		Currency:      &input.Currency,
		Customer:      &input.CustomerID,
		PaymentMethod: &input.PaymentMethodID,
		Confirm:       new(true),
		OffSession:    new(true),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Metadata: input.Metadata,
	}
	createParams.SetIdempotencyKey(input.IdempotencyKey)

	pi, err := g.sc.V1PaymentIntents.Create(ctx, createParams)
	if err != nil {
		return ConfirmOffSessionOutput{}, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	clientSecret := ""
	if pi.Status == stripe.PaymentIntentStatusRequiresAction {
		clientSecret = pi.ClientSecret
	}

	return ConfirmOffSessionOutput{
		PaymentIntentID: pi.ID,
		Status:          string(pi.Status),
		ClientSecret:    clientSecret,
	}, nil
}
