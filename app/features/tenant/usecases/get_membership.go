package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (c *TenantCase) GetMembership(ctx context.Context, userID, tenantID string) (string, error) {
	out, err := c.gateway.GetMember(ctx, g.GetMemberFilter{
		UserID:   &userID,
		TenantID: &tenantID,
	})
	if err != nil {
		return "", err
	}

	if out == nil {
		return "", errs.ErrForbidden
	}

	return out.Role, nil
}
