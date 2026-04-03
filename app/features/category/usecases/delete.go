package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
)

func (c *CategoryCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
