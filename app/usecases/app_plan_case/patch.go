package app_plan_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
)

func (c *AppPlanCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	return c.gateway.Patch(filter, values)
}
