package payment_gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *StripeGatewayAdapter) GetPaymentMethod(ctx context.Context, id string) (*PaymentMethodOutput, error) {
	pm, err := g.sc.V1PaymentMethods.Retrieve(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	if pm.Card == nil {
		return nil, nil
	}

	return &PaymentMethodOutput{
		Brand:    string(pm.Card.Brand),
		Last4:    pm.Card.Last4,
		ExpMonth: pm.Card.ExpMonth,
		ExpYear:  pm.Card.ExpYear,
	}, nil
}
