package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
)

type PatchFilter struct {
	ID *string
}

type PatchValues struct {
	Name string
}

func (c *TenantCase) Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error) {
	return c.gateway.Patch(
		g.PatchFilter{ID: filter.ID},
		g.PatchValues{Name: values.Name},
	)
}
