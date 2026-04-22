package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) GetConfig(ctx context.Context, tenantID string) (*g.GetConfigOutput, error) {
	return s.gateway.GetConfig(ctx, tenantID)
}
