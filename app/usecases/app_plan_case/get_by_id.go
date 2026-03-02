package app_plan_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
)

func (c *AppPlanCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	filter := g.GetByIDFilter{
		ID: id,
	}

	return c.gateway.GetByID(filter)
}
