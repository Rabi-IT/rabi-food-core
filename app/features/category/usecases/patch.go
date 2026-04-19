package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *CategoryCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	logger.GetWideEvent(ctx).SetCategoryID(filter.ID)
	return c.gateway.Patch(ctx, filter, values)
}
