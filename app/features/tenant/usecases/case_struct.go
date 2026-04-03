package usecases

import (
	g "rabi-food-core/features/tenant/gateway"
	user_case "rabi-food-core/features/user/usecases"
)

type TenantCase struct {
	gateway  g.TenantGateway
	userCase *user_case.UserCase
}

func New(
	gateway g.TenantGateway,
	userCase *user_case.UserCase,
) *TenantCase {
	return &TenantCase{
		gateway:  gateway,
		userCase: userCase,
	}
}
