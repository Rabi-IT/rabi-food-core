package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *OrderCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	logger.GetWideEvent(ctx).SetOrderID(filter.ID)

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
