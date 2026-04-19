package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *OrderCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	logger.GetWideEvent(ctx).SetOrderID(filter.ID)
	return c.gateway.Delete(ctx, filter)
}
