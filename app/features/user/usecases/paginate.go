package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
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
