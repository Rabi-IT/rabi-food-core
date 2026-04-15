package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

func (c *TenantCase) GetByID(ctx context.Context, id string, tenantID string) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() || id == "me" {
		id = tenantID
	}

	return c.gateway.GetByID(id)
}
