package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *StripePaymentGatewayAdapter) GetProfile(ctx context.Context, userID string) (*GetProfileOutput, error) {
	user, err := g.auth.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to look up customer: %w", err)
	}

	if user == nil || user.StripePaymentMethodID == "" {
		return &GetProfileOutput{HasPaymentMethod: false}, nil
	}

	pm, err := g.sc.V1PaymentMethods.Retrieve(ctx, user.StripePaymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	if pm.Card == nil {
		return &GetProfileOutput{HasPaymentMethod: true}, nil
	}

	return &GetProfileOutput{
		HasPaymentMethod: true,
		Brand:            string(pm.Card.Brand),
		Last4:            pm.Card.Last4,
		ExpMonth:         pm.Card.ExpMonth,
		ExpYear:          pm.Card.ExpYear,
	}, nil
}
