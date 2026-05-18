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
	profile, err := c.payment.GetPaymentProfile(ctx, userID)
	if err != nil {
		return GetProfileOutput{}, err
	}

	if !profile.HasPaymentMethod {
		return GetProfileOutput{HasPaymentMethod: false}, nil
	}

	if profile.Brand == "" {
		return GetProfileOutput{HasPaymentMethod: true}, nil
	}

	return GetProfileOutput{
		HasPaymentMethod: true,
		Card: &SavedCard{
			Brand:    profile.Brand,
			Last4:    profile.Last4,
			ExpMonth: profile.ExpMonth,
			ExpYear:  profile.ExpYear,
		},
	}, nil
}
