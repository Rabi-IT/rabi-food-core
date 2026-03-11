package subscription_case

import (
	"context"
	"fmt"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/payment"
	"rabi-food-core/domain/subscription"
	"rabi-food-core/libs/database/gateways/product_gateway"
	g "rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"
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

func (s *SubscriptionCase) Create(ctx context.Context, input CreateInput) (string, error) {
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
		return s.handleMissingProducts(ctx, input.Items, products)
	}

	config, err := s.gateway.GetConfig(ctx, session.TenantID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch subscription config: %w", err)
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
		PaymentStatus:       payment.StatusPending,
	})

	if err != nil {
		return "", err
	}

	logger.Get(ctx).Info().
		Str(logger.TenantID, session.TenantID).
		Str(logger.UserID, session.UserID).
		Str(logger.SubscriptionID, id).
		Msg("subscription created")

	return id, nil
}

// handleMissingProducts returns a structured error after identifying which products are absent.
func (s *SubscriptionCase) handleMissingProducts(
	ctx context.Context,
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

	logger.Get(ctx).Warn().
		Strs(logger.ProductID, missingProductIDs).
		Msg("missing products for subscription")

	return "", errs.ProductNotFound(missingProductIDs...)
}
