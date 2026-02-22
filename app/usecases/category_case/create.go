package category_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/category_gateway"
	"rabi-food-core/libs/logger"
)

func (c *CategoryCase) Create(ctx context.Context, input g.CreateInput) (string, error) {
	session := app_context.GetSession(ctx)
	input.TenantID = session.TenantID

	id, err := c.gateway.Create(input)

	if err != nil {
		return "", err
	}

	logger.Get(ctx).Info().Str(logger.TenantID, session.TenantID).Str(logger.CategoryID, id).Msg("category created")

	return id, nil
}
