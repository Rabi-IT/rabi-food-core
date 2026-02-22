package user_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/user_gateway"
)

func (c *UserCase) Patch(ctx context.Context, id string, values g.PatchValues) (bool, error) {
	filter := g.PatchFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Patch(filter, values)
}
