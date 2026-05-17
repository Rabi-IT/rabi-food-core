package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

func (g *StripeGatewayAdapter) CreateCustomer(ctx context.Context, input CreateStripeCustomerInput) (string, error) {
	params := &stripe.CustomerCreateParams{
		Email: stripe.String(input.Email),
		Name:  stripe.String(input.Name),
		Metadata: map[string]string{
			"user_id": input.UserID,
		},
	}
	params.SetIdempotencyKey(input.IdempotencyKey)

	customer, err := g.sc.V1Customers.Create(ctx, params)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	return customer.ID, nil
}
