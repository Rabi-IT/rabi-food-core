package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) Create(ctx context.Context, name string) (string, error) {
	id, err := c.gateway.Create(ctx, g.CreateInput{Name: name})
	logger.GetWideEvent(ctx).SetTenantID(id)
	return id, err
}
