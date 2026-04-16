package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) CreateMember(ctx context.Context, tenantID, userID, role string) error {
	logger.GetWideEvent(ctx).Event = "create-member"
	return c.gateway.CreateMember(ctx, g.CreateMemberInput{
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
	})
}
