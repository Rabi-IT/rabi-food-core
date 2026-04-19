package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
)

func (c *OrderCase) GetByID(ctx context.Context, filter g.GetByIDFilter) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(filter)
}
