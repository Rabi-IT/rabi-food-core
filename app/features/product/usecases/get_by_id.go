package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
)

func (c *ProductCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(ctx, filter)
}
