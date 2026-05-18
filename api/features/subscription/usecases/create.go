package usecases

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database/filter"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

// DeliveryDayInput describes a recurring weekly delivery window provided by the caller.
type DeliveryDayInput struct {
	Weekday   uint8 `json:"weekday"   validate:"min=0,max=6"`
	StartHour uint8 `json:"startHour" validate:"min=0,max=23"`
	EndHour   uint8 `json:"endHour"   validate:"min=1,max=23"`
}

// SubscriptionItemInput is the caller-facing representation of an item.
type SubscriptionItemInput struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  uint   `json:"quantity"  validate:"required,min=1"`
}

// CreateInput contains all data required to create a new subscription group.
type CreateInput struct {
	TenantID     string                  `json:"-"`
	UserID       string                  `json:"-"`
	CheckoutKey  string                  `json:"checkoutKey"`
	Items        []SubscriptionItemInput `json:"items"        validate:"required,min=1"`
	DeliveryDays []g.DeliveryDay         `json:"deliveryDays" validate:"required,min=1"`
	TotalCycles  uint                    `json:"totalCycles"  validate:"required,min=1"`
	Notes        string                  `json:"notes"`
}

func (c *CreateInput) Validate() error {
	seenProducts := make(map[string]struct{}, len(c.Items))
	for _, item := range c.Items {
		if _, exists := seenProducts[item.ProductID]; exists {
			return errs.DuplicateProduct(item.ProductID)
		}
		seenProducts[item.ProductID] = struct{}{}
	}

	seenWeekdays := make(map[uint8]struct{}, len(c.DeliveryDays))
	for _, day := range c.DeliveryDays {
		if _, exists := seenWeekdays[day.Weekday]; exists {
			return errs.DuplicateDeliveryWeekday(day.Weekday)
		}
		seenWeekdays[day.Weekday] = struct{}{}
	}

	return nil
}

// Create creates subscriptions atomically (one per item) and returns their IDs.
func (s *SubscriptionCase) Create(ctx context.Context, input CreateInput) ([]string, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	productIds := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIds = append(productIds, item.ProductID)
	}

	products, err := s.productCase.List(ctx, product_gateway.ListFilter{
		IDs:      productIds,
		IsActive: filter.True,
		TenantID: input.TenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	if len(products) != len(input.Items) {
		return s.handleMissingProducts(input.Items, products)
	}

	config, err := s.gateway.GetConfig(ctx, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription config: %w", err)
	}

	if config == nil {
		return nil, errs.ErrTenantSubscriptionNotConfigured
	}

	price, err := buildSubscriptionPricing(
		input.TotalCycles,
		input.Items,
		products,
		config.DiscountRules,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate price: %w", err)
	}

	gatewayInputs := make([]g.CreateInput, 0, len(price.Products))
	for _, p := range price.Products {
		gatewayInputs = append(gatewayInputs, g.CreateInput{
			UserID:              input.UserID,
			TenantID:            input.TenantID,
			CheckoutKey:         input.CheckoutKey,
			Status:              subscription.StatusActive,
			ProductID:           p.ProductID,
			Quantity:            uint16(p.Quantity),
			UnitPrice:           p.UnitPrice,
			DeliveryDays:        input.DeliveryDays,
			Notes:               input.Notes,
			TotalCycles:         input.TotalCycles,
			RemainingCycles:     input.TotalCycles,
			CycleDiscount:       p.CycleDiscount,
			CutoffOffsetMinutes: config.CutoffOffsetMinutes,
			AutoRenew:           true,
			MaxAttemptsPerOrder: config.MaxAttemptsPerOrder,
			ItemsTotal:          p.ItemsTotal,
			ItemsDiscount:       p.ItemsDiscount,
			PaymentAmount:       p.PaymentAmount(),
			PaymentStatus:       payment_status.StatusPending,
		})
	}

	return s.gateway.BulkCreate(ctx, gatewayInputs)
}

// handleMissingProducts returns a structured error after identifying which products are absent.
func (s *SubscriptionCase) handleMissingProducts(
	requestedItems []SubscriptionItemInput,
	foundProducts []product_gateway.ListOutput,
) ([]string, error) {
	foundProductIDs := make(map[string]struct{}, len(foundProducts))
	for _, product := range foundProducts {
		foundProductIDs[product.ID] = struct{}{}
	}

	missingProductIDs := make([]string, 0)
	for _, item := range requestedItems {
		if _, exists := foundProductIDs[item.ProductID]; !exists {
			missingProductIDs = append(missingProductIDs, item.ProductID)
		}
	}

	return nil, errs.ProductNotFound(missingProductIDs...)
}
