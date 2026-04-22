package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

func (c *TenantCase) GetMembership(ctx context.Context, userID, tenantID string) (*g.GetMemberOutput, error) {
	return c.gateway.GetMember(ctx, g.GetMemberFilter{
		UserID:   userID,
		TenantID: tenantID,
	})
}
