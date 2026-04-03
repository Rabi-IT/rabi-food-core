package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/product/gateway"
)

func (c *ProductCase) Delete(ctx context.Context, id string) (bool, error) {
	filter := g.DeleteFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
