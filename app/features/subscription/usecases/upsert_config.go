package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (s *SubscriptionCase) UpsertConfig(ctx context.Context, input g.UpsertConfigInput) error {
	logger.GetWideEvent(ctx).SetTenantID(input.TenantID)

	return s.gateway.UpsertConfig(ctx, input.TenantID, input)
}
