package category_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/category_gateway"
)

func (c *CategoryCase) Paginate(
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
