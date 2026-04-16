package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) CreateCustomer(ctx context.Context, tenantID string) error {
	session := app_context.GetSession(ctx)
	logger.GetWideEvent(ctx).Event = "create-customer"

	return c.gateway.CreateCustomer(ctx, tenantID, session.UserID)
}
