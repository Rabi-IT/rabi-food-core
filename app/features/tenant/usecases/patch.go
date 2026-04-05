package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type PatchFilter struct {
	ID *string
}

type PatchValues struct {
	Name string
}

func (c *TenantCase) Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error) {
	logger.GetWideEvent(ctx).SetTenantID(*filter.ID)

	return c.gateway.Patch(
		g.PatchFilter{ID: filter.ID},
		g.PatchValues{Name: values.Name},
	)
}
