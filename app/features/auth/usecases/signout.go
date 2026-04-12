package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
)

func (c *AuthCase) SignOut(ctx context.Context, accessToken string) error {
	_ = app_context.GetSession(ctx)
	return c.gateway.SignOut(ctx, accessToken)
}
