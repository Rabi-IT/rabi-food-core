package subscription_case

import (
	"rabi-food-core/libs/database/gateways/product_gateway"
	g "rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/libs/errs"
)

type buildSubscriptionPricingOutput struct {
	SubscriptionItems []g.SubscriptionItem
	ItemsTotal        uint
	ItemsDiscount     uint
	CycleDiscount     uint
}

func (b buildSubscriptionPricingOutput) PaymentAmount() uint {
	return b.ItemsTotal - b.ItemsDiscount - b.CycleDiscount
}

func buildSubscriptionPricing(
	totalCycles uint,
	items []SubscriptionItemInput,
	products []product_gateway.ListOutput,
	cycleDiscountRules []g.DiscountRule,
) (buildSubscriptionPricingOutput, error) {
	output := buildSubscriptionPricingOutput{
		SubscriptionItems: make([]g.SubscriptionItem, 0, len(items)),
	}

	productMap := make(map[string]product_gateway.ListOutput, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return buildSubscriptionPricingOutput{}, errs.ProductNotFound(item.ProductID)
		}

		subtotal := product.Price * item.Quantity
		discountPercent := uint(0)
		for _, rule := range product.DiscountRules {
			if item.Quantity >= rule.QuantityThreshold && rule.Discount > discountPercent {
				discountPercent = rule.Discount
			}
		}

		output.SubscriptionItems = append(output.SubscriptionItems, g.SubscriptionItem{
			ItemID:    product.ID,
			Name:      product.Name,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
			Discount:  discountPercent,
			Total:     subtotal - subtotal*discountPercent/100,
		})

		output.ItemsTotal += subtotal
		output.ItemsDiscount += subtotal * discountPercent / 100 //nolint:mnd
	}

	cycleDiscountPercent := uint8(0)
	for _, rule := range cycleDiscountRules {
		if totalCycles >= rule.CyclesThreshold && rule.Discount > cycleDiscountPercent {
			cycleDiscountPercent = rule.Discount
		}
	}

	output.CycleDiscount = (output.ItemsTotal - output.ItemsDiscount) * uint(cycleDiscountPercent) / 100 //nolint:mnd

	return output, nil
}
