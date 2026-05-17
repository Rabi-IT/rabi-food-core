package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

func (g *StripeGatewayAdapter) CreatePaymentIntent(ctx context.Context, input CreateIntentInput) (CreateIntentOutput, error) {
	params := &stripe.PaymentIntentCreateParams{
		Amount:   stripe.Int64(input.AmountCents),
		Currency: stripe.String(input.Currency),
		Customer: stripe.String(input.CustomerID),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		SetupFutureUsage: stripe.String("off_session"),
		Metadata:         input.Metadata,
	}
	params.SetIdempotencyKey(input.IdempotencyKey)

	pi, err := g.sc.V1PaymentIntents.Create(ctx, params)
	if err != nil {
		return CreateIntentOutput{}, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	return CreateIntentOutput{
		PaymentIntentID: pi.ID,
		ClientSecret:    pi.ClientSecret,
	}, nil
}
