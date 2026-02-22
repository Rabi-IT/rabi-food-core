package tenant_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
)

func (c *TenantCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		id = session.TenantID
	}

	if id == "me" {
		id = session.TenantID
	}

	return c.gateway.GetByID(id)
}
