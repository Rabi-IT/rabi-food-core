package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (c *OrderCase) Paginate(
	ctx context.Context,
	input g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	return c.gateway.Paginate(ctx, input, paginate)
}
