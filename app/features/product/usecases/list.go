package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/product/gateway"
)

func (c *ProductCase) List(ctx context.Context, filter g.ListFilter) ([]g.ListOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.List(filter)
}
