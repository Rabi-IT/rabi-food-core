package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	return s.gateway.GetByID(ctx, filter)
}
