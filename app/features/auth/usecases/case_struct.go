package usecases

import "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"

type AuthCase struct {
	gateway gateway.AuthGateway
}

func New(gateway gateway.AuthGateway) *AuthCase {
	return &AuthCase{gateway: gateway}
}
