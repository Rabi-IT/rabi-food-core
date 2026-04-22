package usecases

import (
	"github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	tenant_usecases "github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
)

type AuthCase struct {
	gateway    gateway.AuthGateway
	tenantCase *tenant_usecases.TenantCase
}

func New(gateway gateway.AuthGateway, tenantCase *tenant_usecases.TenantCase) *AuthCase {
	return &AuthCase{
		gateway:    gateway,
		tenantCase: tenantCase,
	}
}
