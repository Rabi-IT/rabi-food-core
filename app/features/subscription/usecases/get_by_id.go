package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.UserID = session.UserID
	} else {
		filter.TenantID = ""
	}

	return s.gateway.GetByID(ctx, filter)
}
