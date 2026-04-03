// We want to test `buildSubscriptionPricing` without exporting it from the `subscription_case` package.
//
//nolint:testpackage
package usecases

import (
	product_gateway "rabi-food-core/features/product/gateway"
	subscription_gateway "rabi-food-core/features/subscription/gateway"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type buildSubscriptionPricingTest struct {
	name               string
	totalCycles        uint
	items              []SubscriptionItemInput
	products           []product_gateway.ListOutput
	cycleDiscountRules []subscription_gateway.DiscountRule
	expectedOutput     buildSubscriptionPricingOutput
	expectError        bool
}

func (b *buildSubscriptionPricingTest) run(t *testing.T) {
	t.Helper()
	t.Run(b.name, func(t *testing.T) {
		output, err := buildSubscriptionPricing(b.totalCycles, b.items, b.products, b.cycleDiscountRules)
		if b.expectError {
			require.Error(t, err)

			return
		}

		require.NoError(t, err)
		assert.Equal(t, b.expectedOutput.ItemsTotal, output.ItemsTotal, "ItemsTotal mismatch")
		assert.Equal(t, b.expectedOutput.ItemsDiscount, output.ItemsDiscount, "ItemsDiscount mismatch")
		assert.Equal(t, b.expectedOutput.CycleDiscount, output.CycleDiscount, "CycleDiscount mismatch")
	})
}

func Test_BuildSubscriptionPricing__SingleProduct(t *testing.T) {
	var tests = []buildSubscriptionPricingTest{
		{
			name:        "when no discount rules apply the product subtotal becomes the subscription total",
			totalCycles: 3,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products: []product_gateway.ListOutput{{
				ID:            "prod1",
				Price:         100,
				DiscountRules: nil,
			}},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 200, ItemsDiscount: 0, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when quantity reaches the discount threshold the product discount is applied",
			totalCycles: 1,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 5}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
				},
			}},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 500, ItemsDiscount: 50, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when quantity does not reach the discount threshold the product discount is not applied",
			totalCycles: 1,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 4}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
				},
			}},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 400, ItemsDiscount: 0, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when multiple product discount rules match the highest discount is applied",
			totalCycles: 1,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 10}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
					{QuantityThreshold: 10, Discount: 20},
				},
			}},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 1000, ItemsDiscount: 200, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when cycle discount rules also apply the cycle discount is combined with the product discount",
			totalCycles: 3,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 5}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
				},
			}},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
			},
			expectedOutput: buildSubscriptionPricingOutput{ItemsTotal: 500, ItemsDiscount: 50, CycleDiscount: 22},
			expectError:    false,
		},
		{
			name: "when quantity does not meet the discount threshold the product discount is not applied" +
				"even if a cycle discount rule matches",
			totalCycles: 2,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
				},
			}},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
			},
			expectedOutput: buildSubscriptionPricingOutput{ItemsTotal: 200, ItemsDiscount: 0, CycleDiscount: 10},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		tt.run(t)
	}
}

func Test_BuildSubscriptionPricing__MultipleProducts(t *testing.T) {
	var tests = []buildSubscriptionPricingTest{
		{
			name:        "when no discount rules apply the product subtotal becomes the subscription total",
			totalCycles: 3,
			items: []SubscriptionItemInput{
				{ProductID: "prod1", Quantity: 2},
				{ProductID: "prod2", Quantity: 3},
			},
			products: []product_gateway.ListOutput{
				{ID: "prod1", Price: 100, DiscountRules: nil},
				{ID: "prod2", Price: 50, DiscountRules: nil},
			},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 350, ItemsDiscount: 0, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when products have different discount rules each product applies its own discount independently",
			totalCycles: 1,
			items: []SubscriptionItemInput{
				{ProductID: "prod1", Quantity: 5},
				{ProductID: "prod2", Quantity: 10},
			},
			products: []product_gateway.ListOutput{
				{
					ID:    "prod1",
					Price: 100,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 5, Discount: 10},
					},
				},
				{
					ID:    "prod2",
					Price: 50,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 10, Discount: 20},
					},
				},
			},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 1000, ItemsDiscount: 150, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when one product qualifies for a discount and another does not each product is priced independently",
			totalCycles: 1,
			items: []SubscriptionItemInput{
				{ProductID: "prod1", Quantity: 5},
				{ProductID: "prod2", Quantity: 9},
			},
			products: []product_gateway.ListOutput{
				{
					ID:    "prod1",
					Price: 100,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 5, Discount: 10},
					},
				},
				{
					ID:    "prod2",
					Price: 50,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 10, Discount: 20},
					},
				},
			},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 950, ItemsDiscount: 50, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when only some products have discount rules only those products receive discounts",
			totalCycles: 1,
			items: []SubscriptionItemInput{
				{ProductID: "prod1", Quantity: 5},
				{ProductID: "prod2", Quantity: 10},
			},
			products: []product_gateway.ListOutput{
				{
					ID:    "prod1",
					Price: 100,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 5, Discount: 10},
					},
				},
				{
					ID:            "prod2",
					Price:         50,
					DiscountRules: nil,
				},
			},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 1000, ItemsDiscount: 50, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name: "when multiple product discount rules match then both highest discounts are applied independently" +
				"for each product",
			totalCycles: 1,
			items: []SubscriptionItemInput{
				{
					ProductID: "prod1", Quantity: 10,
				},
				{
					ProductID: "prod2", Quantity: 15,
				},
			},
			products: []product_gateway.ListOutput{
				{
					ID:    "prod1",
					Price: 100,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 5, Discount: 10},
						{QuantityThreshold: 10, Discount: 20},
					},
				},
				{
					ID:    "prod2",
					Price: 50,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 5, Discount: 5},
						{QuantityThreshold: 10, Discount: 10},
					},
				},
			},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{ItemsTotal: 1750, ItemsDiscount: 275, CycleDiscount: 0},
			expectError:        false,
		},
		{
			name:        "when multiple product and cycle discount rules match then both highest discounts are applied",
			totalCycles: 6,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 10}},
			products: []product_gateway.ListOutput{{
				ID:    "prod1",
				Price: 100,
				DiscountRules: []product_gateway.DiscountRule{
					{QuantityThreshold: 5, Discount: 10},
					{QuantityThreshold: 10, Discount: 20},
				},
			}},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
				{CyclesThreshold: 5, Discount: 10},
			},
			expectedOutput: buildSubscriptionPricingOutput{ItemsTotal: 1000, ItemsDiscount: 200, CycleDiscount: 80},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		tt.run(t)
	}
}

func Test_BuildSubscriptionPricing__CycleDiscount(t *testing.T) {
	var tests = []buildSubscriptionPricingTest{
		{
			name:        "when the subscription cycles reach the cycle discount threshold the cycle discount is applied",
			totalCycles: 3,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products:    []product_gateway.ListOutput{{ID: "prod1", Price: 100}},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
			},
			expectedOutput: buildSubscriptionPricingOutput{ItemsTotal: 200, ItemsDiscount: 0, CycleDiscount: 10},
			expectError:    false,
		},
		{
			name:        "when multiple cycle discount rules match the highest discount is applied",
			totalCycles: 6,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products:    []product_gateway.ListOutput{{ID: "prod1", Price: 100}},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
				{CyclesThreshold: 5, Discount: 10},
			},
			expectedOutput: buildSubscriptionPricingOutput{ItemsTotal: 200, ItemsDiscount: 0, CycleDiscount: 20},
			expectError:    false,
		},
		{
			name:        "when cycle discount rules exist but no threshold is reached no cycle discount is applied",
			totalCycles: 1,
			items: []SubscriptionItemInput{
				{ProductID: "prod1", Quantity: 2},
			},
			products: []product_gateway.ListOutput{
				{ID: "prod1", Price: 100},
			},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 2, Discount: 5},
			},
			expectedOutput: buildSubscriptionPricingOutput{
				ItemsTotal:    200,
				ItemsDiscount: 0,
				CycleDiscount: 0,
			},
			expectError: false,
		},
		{
			name: "when cycle discount rules exist but no threshold is reached no cycle discount is applied " +
				"even if product quantity matches product discount rule",
			totalCycles: 1,
			items:       []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products: []product_gateway.ListOutput{
				{
					ID:    "prod1",
					Price: 100,
					DiscountRules: []product_gateway.DiscountRule{
						{QuantityThreshold: 2, Discount: 10},
					},
				},
			},
			cycleDiscountRules: []subscription_gateway.DiscountRule{
				{CyclesThreshold: 5, Discount: 10},
			},
			expectedOutput: buildSubscriptionPricingOutput{
				ItemsTotal:    200,
				ItemsDiscount: 20,
				CycleDiscount: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt.run(t)
	}
}

func Test_BuildSubscriptionPricing__Errors(t *testing.T) {
	var tests = []buildSubscriptionPricingTest{
		{
			name:               "should return error if product not found",
			totalCycles:        1,
			items:              []SubscriptionItemInput{{ProductID: "prod1", Quantity: 2}},
			products:           []product_gateway.ListOutput{},
			cycleDiscountRules: nil,
			expectedOutput:     buildSubscriptionPricingOutput{},
			expectError:        true,
		},
	}

	for _, tt := range tests {
		tt.run(t)
	}
}
