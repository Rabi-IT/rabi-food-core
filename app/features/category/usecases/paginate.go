package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (c *CategoryCase) Paginate(
	ctx context.Context,
	input g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	return c.gateway.Paginate(input, paginate)
}
