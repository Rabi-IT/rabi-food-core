package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"
)

func (c *CategoryCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(filter)
}
