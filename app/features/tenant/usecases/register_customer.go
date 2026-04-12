package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) RegisterCustomer(ctx context.Context, tenantID string) error {
	session := app_context.GetSession(ctx)
	logger.GetWideEvent(ctx).Event = "register-customer"

	return c.gateway.RegisterCustomer(ctx, tenantID, session.UserID)
}
