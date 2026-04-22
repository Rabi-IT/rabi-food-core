package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) CreateMember(ctx context.Context, tenantID, userID string, role tenant.Role) error {
	logger.GetWideEvent(ctx).Event = "create-member"
	return c.gateway.CreateMember(ctx, g.CreateMemberInput{
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
	})
}
