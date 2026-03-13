package subscription_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/subscription_gateway"
)

func (s *SubscriptionCase) UpsertConfig(ctx context.Context, input g.UpsertConfigInput) (bool, error) {
	session := app_context.GetSession(ctx)

	upserted, err := s.gateway.UpsertConfig(ctx, session.TenantID, input)

	if err != nil {
		return false, err
	}

	return upserted, nil
}
