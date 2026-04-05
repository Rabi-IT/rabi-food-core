package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *ProductCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	logger.GetWideEvent(ctx).SetProductID(filter.ID)

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.Patch(filter, values)
}
