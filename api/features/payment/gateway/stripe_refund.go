package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

func (g *StripePaymentGatewayAdapter) Refund(ctx context.Context, paymentIntentID string) error {
	params := &stripe.RefundCreateParams{
		PaymentIntent: &paymentIntentID,
	}

	_, err := g.sc.V1Refunds.Create(ctx, params)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	return nil
}
