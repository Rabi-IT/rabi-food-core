package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
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

	logger.GetWideEvent(ctx).ProductID = id

	return id, nil
}
