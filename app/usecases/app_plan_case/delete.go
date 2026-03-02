package app_plan_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
)

func (c *AppPlanCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {

	return c.gateway.Delete(filter)
}
