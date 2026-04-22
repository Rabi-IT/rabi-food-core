package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) EnrollCustomer(ctx context.Context, tenantID string) error {
	session := app_context.GetSession(ctx)
	logger.GetWideEvent(ctx).Event = "enroll-customer"

	return c.gateway.CreateCustomer(ctx, tenantID, session.UserID)
}
