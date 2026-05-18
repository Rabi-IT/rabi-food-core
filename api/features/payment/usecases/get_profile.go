package payment_usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
)

func (c *PaymentCase) GetProfile(ctx context.Context, userID string) (*g.GetProfileOutput, error) {
	return c.payment.GetProfile(ctx, userID)
}
