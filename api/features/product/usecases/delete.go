package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *ProductCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	logger.GetWideEvent(ctx).SetProductID(filter.ID)

	return c.gateway.Delete(ctx, filter)
}
