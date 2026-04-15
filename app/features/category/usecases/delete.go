package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *CategoryCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	logger.GetWideEvent(ctx).SetCategoryID(filter.ID)
	return c.gateway.Delete(filter)
}
