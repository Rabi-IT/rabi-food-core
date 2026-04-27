package usecases

import (
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type productPricing struct {
	ProductID     string
	Quantity      uint
	UnitPrice     uint
	ItemsTotal    uint
	ItemsDiscount uint
	CycleDiscount uint
}

func (p productPricing) PaymentAmount() uint {
	return p.ItemsTotal - p.ItemsDiscount - p.CycleDiscount
}

type buildSubscriptionPricingOutput struct {
	Products      []productPricing
	ItemsTotal    uint
	ItemsDiscount uint
	CycleDiscount uint
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
		Products: make([]productPricing, 0, len(items)),
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
		itemDiscount := subtotal * discountPercent / 100 //nolint:mnd

		output.Products = append(output.Products, productPricing{
			ProductID:     product.ID,
			Quantity:      item.Quantity,
			UnitPrice:     product.Price,
			ItemsTotal:    subtotal,
			ItemsDiscount: itemDiscount,
		})

		output.ItemsTotal += subtotal
		output.ItemsDiscount += itemDiscount
	}

	cycleDiscountPercent := uint8(0)
	for _, rule := range cycleDiscountRules {
		if totalCycles >= rule.CyclesThreshold && rule.Discount > cycleDiscountPercent {
			cycleDiscountPercent = rule.Discount
		}
	}

	output.CycleDiscount = (output.ItemsTotal - output.ItemsDiscount) * uint(cycleDiscountPercent) / 100 //nolint:mnd

	for i := range output.Products {
		net := output.Products[i].ItemsTotal - output.Products[i].ItemsDiscount
		output.Products[i].CycleDiscount = net * uint(cycleDiscountPercent) / 100 //nolint:mnd
	}

	return output, nil
}
