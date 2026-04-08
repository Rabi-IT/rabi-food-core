package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
)

func (c *ProductCase) List(ctx context.Context, filter g.ListFilter) ([]g.ListOutput, error) {
	return c.gateway.List(filter)
}
