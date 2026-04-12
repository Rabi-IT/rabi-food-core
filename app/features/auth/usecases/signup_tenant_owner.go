package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

// SignUpTenantOwner satisfies the tenant.userCreator interface.
// It creates the initial owner account when a new tenant is registered.
func (c *AuthCase) SignUpTenantOwner(ctx context.Context, email, password, name, phone, tenantID string) (string, error) {
	out, err := c.gateway.SignUp(ctx, g.SignUpInput{
		Email:    email,
		Password: password,
		Name:     name,
		Phone:    phone,
		TenantID: tenantID,
		Role:     string(auth.User),
	})
	if err != nil {
		return "", err
	}

	return out.ID, nil
}
