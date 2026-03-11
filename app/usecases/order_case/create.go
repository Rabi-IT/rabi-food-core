package order_case

import (
	"context"
	"fmt"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/order"
	"rabi-food-core/domain/payment"
	g "rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"

	"github.com/google/uuid"
)

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  uint   `json:"quantity"`
}

type CreateInput struct {
	Items []OrderItem `json:"items"`
	Notes string      `json:"notes"`
}

func (c *OrderCase) Create(ctx context.Context, input CreateInput) (string, error) {
	productIds := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIds = append(productIds, item.ProductID)
	}

	isActive := true
	products, err := c.productCase.List(ctx, product_gateway.ListFilter{
		IDs:      productIds,
		IsActive: &isActive,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch products: %w", err)
	}

	if len(products) != len(input.Items) {
		return c.handleMissingProducts(ctx, input.Items, products)
	}

	productMap := make(map[string]product_gateway.ListOutput)
	for _, product := range products {
		productMap[product.ID] = product
	}

	orderItems := make([]g.OrderItem, 0, len(input.Items))
	totalPrice := uint(0)
	for _, item := range input.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			logger.Get(ctx).Warn().Msgf("product not found for ID: %s", item.ProductID)

			return "", errs.ProductNotFound(item.ProductID)
		}

		itemTotal := product.Price * item.Quantity
		orderItems = append(orderItems, g.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			Total:       itemTotal,
		})
		totalPrice += itemTotal
	}

	session := app_context.GetSession(ctx)
	id, err := c.gateway.Create(g.CreateInput{
		UserID:   session.UserID,
		TenantID: session.TenantID,

		Code:              uuid.Must(uuid.NewV7()).String(),
		PaymentStatus:     payment.StatusPending,
		FulfillmentStatus: order.FulfillmentPending,
		DeliveryStatus:    order.DeliveryPending,
		Notes:             input.Notes,
		TotalPrice:        totalPrice,
		Items:             orderItems,
	})

	if err != nil {
		return "", err
	}

	logger.Get(ctx).Info().Str(logger.TenantID, session.TenantID).Str(logger.OrderID, id).Msg("order created")

	return id, nil
}

func (c *OrderCase) handleMissingProducts(
	ctx context.Context,
	requestedItems []OrderItem,
	foundProducts []product_gateway.ListOutput,
) (string, error) {
	foundProductIDs := make(map[string]struct{})
	for _, product := range foundProducts {
		foundProductIDs[product.ID] = struct{}{}
	}

	missingProductIDs := make([]string, 0)
	for _, item := range requestedItems {
		if _, exists := foundProductIDs[item.ProductID]; !exists {
			missingProductIDs = append(missingProductIDs, item.ProductID)
		}
	}

	logger.Get(ctx).Warn().Msgf("missing products: %v", missingProductIDs)

	return "", errs.ProductNotFound(missingProductIDs...)
}
