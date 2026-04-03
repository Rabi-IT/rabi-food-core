package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) UpsertConfig(ctx context.Context, input g.UpsertConfigInput) error {
	session := app_context.GetSession(ctx)

	return s.gateway.UpsertConfig(ctx, session.TenantID, input)
}
