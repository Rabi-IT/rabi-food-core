package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

func (c *TenantCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(id)
}
