package order_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/order_gateway"
)

func (c *OrderCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
