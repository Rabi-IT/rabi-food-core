package product_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/logger"
)

func (c *ProductCase) Create(ctx context.Context, input g.CreateInput) (string, error) {
	session := app_context.GetSession(ctx)
	if !session.Role.IsBackoffice() {
		input.TenantID = session.TenantID
	}

	id, err := c.gateway.Create(input)

	if err != nil {
		return "", err
	}

	logger.Get(ctx).Info().Str(logger.TenantID, session.TenantID).Str(logger.ProductID, id).Msg("product created")

	return id, nil
}
