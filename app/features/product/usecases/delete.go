package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *ProductCase) Delete(ctx context.Context, id string) (bool, error) {
	logger.GetWideEvent(ctx).SetProductID(id)

	filter := g.DeleteFilter{
		ID: id,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Delete(filter)
}
