package payment_usecases

import "context"

type SavedCard struct {
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	ExpMonth int64  `json:"expMonth"`
	ExpYear  int64  `json:"expYear"`
}

type GetProfileOutput struct {
	HasPaymentMethod bool       `json:"hasPaymentMethod"`
	Card             *SavedCard `json:"card,omitempty"`
}

func (c *PaymentCase) GetProfile(ctx context.Context, userID string) (GetProfileOutput, error) {
	record, err := c.customer.GetByUserID(ctx, userID)
	if err != nil {
		return GetProfileOutput{}, err
	}

	if record == nil || !record.HasPaymentMethod || record.StripePaymentMethodID == "" {
		return GetProfileOutput{HasPaymentMethod: false}, nil
	}

	pm, err := c.stripe.GetPaymentMethod(ctx, record.StripePaymentMethodID)
	if err != nil {
		return GetProfileOutput{}, err
	}

	if pm == nil {
		return GetProfileOutput{HasPaymentMethod: true}, nil
	}

	return GetProfileOutput{
		HasPaymentMethod: true,
		Card: &SavedCard{
			Brand:    pm.Brand,
			Last4:    pm.Last4,
			ExpMonth: pm.ExpMonth,
			ExpYear:  pm.ExpYear,
		},
	}, nil
}
