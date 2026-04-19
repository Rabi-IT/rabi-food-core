package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *ProductCase) Create(ctx context.Context, input g.CreateInput) (string, error) {
	id, err := c.gateway.Create(ctx, input)
	logger.GetWideEvent(ctx).SetProductID(id)
	if err != nil {
		return "", err
	}

	return id, nil
}
