package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
)

func (c *CategoryCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	filter := g.GetByIDFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.GetByID(filter)
}
