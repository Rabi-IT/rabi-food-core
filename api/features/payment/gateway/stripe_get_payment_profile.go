package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *StripePaymentGatewayAdapter) GetPaymentProfile(ctx context.Context, userID string) (*PaymentProfileOutput, error) {
	meta, err := g.getUserMeta(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to look up customer: %w", err)
	}

	if meta.StripePaymentMethodID == "" {
		return &PaymentProfileOutput{HasPaymentMethod: false}, nil
	}

	pm, err := g.sc.V1PaymentMethods.Retrieve(ctx, meta.StripePaymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	if pm.Card == nil {
		return &PaymentProfileOutput{HasPaymentMethod: true}, nil
	}

	return &PaymentProfileOutput{
		HasPaymentMethod: true,
		Brand:            string(pm.Card.Brand),
		Last4:            pm.Card.Last4,
		ExpMonth:         pm.Card.ExpMonth,
		ExpYear:          pm.Card.ExpYear,
	}, nil
}
