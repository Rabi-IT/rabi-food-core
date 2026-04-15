package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
)

func (c *OrderCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	if !session.Role.IsUser() {
		filter.TenantID = ""
	}

	return c.gateway.GetByID(filter)
}
