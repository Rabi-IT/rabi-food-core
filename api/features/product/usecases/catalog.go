package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database/filter"
)

const maxProductsPerTenant uint64 = 50

type CatalogInput struct {
	TenantID string
}

func (c *ProductCase) Catalog(ctx context.Context, input CatalogInput) ([]g.ListOutput, error) {
	return c.gateway.List(ctx, g.ListFilter{
		TenantID:         input.TenantID,
		IsActive:         filter.True,
		WithCategoryName: true,
		Limit:            maxProductsPerTenant,
	})
}
