package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (c *OrderCase) Paginate(
	ctx context.Context,
	input g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		input.TenantID = &session.TenantID
	}

	return c.gateway.Paginate(input, paginate)
}
