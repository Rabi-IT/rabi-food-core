package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
)

func (c *UserCase) Patch(ctx context.Context, id string, values g.PatchValues) (bool, error) {
	filter := g.PatchFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Patch(filter, values)
}
