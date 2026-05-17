package usecases

import (
	"context"
	"fmt"

	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database/filter"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

// CalculateTotalAmount returns the total amount in cents for the full subscription plan
// (per-cycle payment amount × number of cycles), applying all applicable discounts.
func (s *SubscriptionCase) CalculateTotalAmount(ctx context.Context, tenantID string, items []SubscriptionItemInput, totalCycles uint) (uint, error) {
	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productCase.List(ctx, product_gateway.ListFilter{
		IDs:      productIDs,
		IsActive: filter.True,
		TenantID: tenantID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	config, err := s.gateway.GetConfig(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch subscription config: %w", err)
	}

	if config == nil {
		return 0, errs.ErrTenantSubscriptionNotConfigured
	}

	pricing, err := buildSubscriptionPricing(totalCycles, items, products, config.DiscountRules)
	if err != nil {
		return 0, err
	}

	return pricing.PaymentAmount() * totalCycles, nil
}
