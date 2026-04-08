package usecases

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/order"
	g "github.com/Rabi-IT/rabi-food-core/features/order/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/google/uuid"
)

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  uint   `json:"quantity"`
}

type CreateInput struct {
	TenantID string      `json:"-"`
	UserID   string      `json:"-"`
	Items    []OrderItem `json:"items"`
	Notes    string      `json:"notes"`
}

func (c *OrderCase) Create(ctx context.Context, input CreateInput) (string, error) {
	productIds := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIds = append(productIds, item.ProductID)
	}

	products, err := c.productCase.List(ctx, product_gateway.ListFilter{
		IDs:      productIds,
		IsActive: new(true),
		TenantID: &input.TenantID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch products: %w", err)
	}

	if len(products) != len(input.Items) {
		return c.handleMissingProducts(input.Items, products)
	}

	productMap := make(map[string]gateway.ListOutput)
	for _, product := range products {
		productMap[product.ID] = product
	}

	orderItems := make([]g.OrderItem, 0, len(input.Items))
	totalPrice := uint(0)
	for _, item := range input.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
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

	id, err := c.gateway.Create(g.CreateInput{
		UserID:   input.UserID,
		TenantID: input.TenantID,

		Code:              uuid.Must(uuid.NewV7()).String(),
		PaymentStatus:     payment_status.StatusPending,
		FulfillmentStatus: order.FulfillmentPending,
		DeliveryStatus:    order.DeliveryPending,
		Notes:             input.Notes,
		TotalPrice:        totalPrice,
		Items:             orderItems,
	})

	logger.GetWideEvent(ctx).SetOrderID(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *OrderCase) handleMissingProducts(
	requestedItems []OrderItem,
	foundProducts []gateway.ListOutput,
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

	return "", errs.ProductNotFound(missingProductIDs...)
}
