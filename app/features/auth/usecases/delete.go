package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
)

func (c *AuthCase) Delete(ctx context.Context, id string) error {
	session := app_context.GetSession(ctx)

	// Users can only delete themselves
	if session.Role.IsUser() {
		id = session.UserID
	}

	return c.gateway.Delete(ctx, id)
}
