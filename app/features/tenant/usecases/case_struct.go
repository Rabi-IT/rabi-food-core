package usecases

import (
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

type TenantCase struct {
	gateway g.TenantGateway
}

func New(gateway g.TenantGateway) *TenantCase {
	return &TenantCase{gateway: gateway}
}
