package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (s *SubscriptionCase) UpsertConfig(ctx context.Context, input g.UpsertConfigInput) error {
	session := app_context.GetSession(ctx)

	logger.GetWideEvent(ctx).SetTenantID(session.TenantID)

	return s.gateway.UpsertConfig(ctx, session.TenantID, input)
}
