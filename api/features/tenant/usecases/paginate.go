package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (c *TenantCase) Paginate(
	ctx context.Context,
	filter g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	return c.gateway.Paginate(ctx, filter, paginate)
}
