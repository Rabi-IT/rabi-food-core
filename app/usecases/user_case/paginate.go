package user_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/user_gateway"
)

func (c *UserCase) Paginate(
	ctx context.Context,
	filter g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = &session.TenantID
	}

	return c.gateway.Paginate(filter, paginate)
}
