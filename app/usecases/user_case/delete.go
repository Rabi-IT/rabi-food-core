package user_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/libs/database/gateways/user_gateway"
)

func (c *UserCase) Delete(ctx context.Context, id string) (bool, error) {
	filter := user_gateway.DeleteFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
