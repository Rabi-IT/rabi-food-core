package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *StripeGatewayAdapter) GetPaymentIntent(ctx context.Context, id string) (*PaymentIntentOutput, error) {
	pi, err := g.sc.V1PaymentIntents.Retrieve(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	paymentMethodID := ""
	if pi.PaymentMethod != nil {
		paymentMethodID = pi.PaymentMethod.ID
	}

	return &PaymentIntentOutput{
		ID:              pi.ID,
		Status:          string(pi.Status),
		PaymentMethodID: paymentMethodID,
	}, nil
}
