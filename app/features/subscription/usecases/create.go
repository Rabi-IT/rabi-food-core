package usecases

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
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

// CreateInput contains all data required to create a new subscription.
type CreateInput struct {
	Items        []SubscriptionItemInput `json:"items"        validate:"required,min=1"`
	DeliveryDays []g.DeliveryDay         `json:"deliveryDays" validate:"required,min=1"`
	TotalCycles  uint                    `json:"totalCycles"  validate:"required,min=1"`
	AutoRenew    bool                    `json:"autoRenew"`
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

func (s *SubscriptionCase) Create(ctx context.Context, input CreateInput) (string, error) {
	err := input.Validate()
	if err != nil {
		return "", err
	}

	productIds := make([]string, 0, len(input.Items))
	session := app_context.GetSession(ctx)

	for _, item := range input.Items {
		productIds = append(productIds, item.ProductID)
	}

	products, err := s.productCase.List(ctx, product_gateway.ListFilter{
		IDs:      productIds,
		IsActive: new(true),
		TenantID: session.TenantID,
	})

	if err != nil {
		return "", fmt.Errorf("failed to fetch products: %w", err)
	}

	if len(products) != len(input.Items) {
		return s.handleMissingProducts(input.Items, products)
	}

	config, err := s.gateway.GetConfig(ctx, session.TenantID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch subscription config: %w", err)
	}

	if config == nil {
		return "", errs.ErrTenantSubscriptionNotConfigured
	}

	price, err := buildSubscriptionPricing(
		input.TotalCycles,
		input.Items,
		products,
		config.DiscountRules,
	)

	if err != nil {
		return "", fmt.Errorf("failed to calculate price: %w", err)
	}

	id, err := s.gateway.Create(ctx, g.CreateInput{
		UserID:              session.UserID,
		TenantID:            session.TenantID,
		Status:              subscription.StatusActive,
		DeliveryDays:        input.DeliveryDays,
		Items:               price.SubscriptionItems,
		Notes:               input.Notes,
		TotalCycles:         input.TotalCycles,
		RemainingCycles:     input.TotalCycles,
		CycleDiscount:       price.CycleDiscount,
		CutoffOffsetMinutes: config.CutoffOffsetMinutes,
		AutoRenew:           input.AutoRenew,
		MaxAttemptsPerOrder: config.MaxAttemptsPerOrder,
		ItemsTotal:          price.ItemsTotal,
		ItemsDiscount:       price.ItemsDiscount,
		PaymentAmount:       price.PaymentAmount(),
		PaymentStatus:       payment_status.StatusPending,
	})
	logger.GetWideEvent(ctx).SetSubscriptionID(id)

	if err != nil {
		return "", err
	}

	return id, nil
}

// handleMissingProducts returns a structured error after identifying which products are absent.
func (s *SubscriptionCase) handleMissingProducts(
	requestedItems []SubscriptionItemInput,
	foundProducts []product_gateway.ListOutput,
) (string, error) {
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

	return "", errs.ProductNotFound(missingProductIDs...)
}
