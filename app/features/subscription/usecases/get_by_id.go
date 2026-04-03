package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	filter := g.GetByIDFilter{ID: id}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
		filter.UserID = session.UserID
	}

	return s.gateway.GetByID(ctx, filter)
}
