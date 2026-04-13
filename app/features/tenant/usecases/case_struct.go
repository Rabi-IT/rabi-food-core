package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

// userCreator is satisfied by auth_case.AuthCase — defined here to avoid cross-feature imports.
// Primitive parameters are used intentionally so auth does not need to import tenant types.
type userCreator interface {
	SignUpTenantOwner(ctx context.Context, email, password, name, phone string) (string, error)
}

type TenantCase struct {
	gateway     g.TenantGateway
	userCreator userCreator
}

func New(gateway g.TenantGateway, userCreator userCreator) *TenantCase {
	return &TenantCase{
		gateway:     gateway,
		userCreator: userCreator,
	}
}
