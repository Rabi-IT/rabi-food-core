package usecases

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/features/product/gateway"
	"rabi-food-core/libs/database"
)

func (c *ProductCase) Paginate(
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
