package app_plan_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
	"rabi-food-core/libs/logger"
)

func (c *AppPlanCase) Create(ctx context.Context, input g.CreateInput) (string, error) {
	id, err := c.gateway.Create(input)

	if err != nil {
		return "", err
	}

	logger.Get(ctx).Info().Str(logger.AppPlanID, id).Msg("category created")

	return id, nil
}
