package app_plan_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
)

func (c *AppPlanCase) List(ctx context.Context, filter g.ListFilter) ([]g.ListOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.List(filter)
}
