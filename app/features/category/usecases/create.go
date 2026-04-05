package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *CategoryCase) Create(ctx context.Context, input g.CreateInput) (string, error) {
	session := app_context.GetSession(ctx)
	input.TenantID = session.TenantID

	id, err := c.gateway.Create(input)
	logger.GetWideEvent(ctx).SetCategoryID(id)
	if err != nil {
		return "", err
	}

	return id, nil
}
