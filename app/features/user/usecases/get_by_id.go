package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/user/gateway"
)

func (c *UserCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	if filter.ID == "me" {
		filter.ID = ""
	}

	return c.gateway.GetByID(filter)
}
