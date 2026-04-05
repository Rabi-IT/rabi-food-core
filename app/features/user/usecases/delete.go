package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *UserCase) Delete(ctx context.Context, id string) (bool, error) {
	filter := g.DeleteFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	logger.GetWideEvent(ctx).SetUserID(id)
	return c.gateway.Delete(filter)
}
