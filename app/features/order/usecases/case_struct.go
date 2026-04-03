package usecases

import (
	g "rabi-food-core/features/order/gateway"
	pc "rabi-food-core/features/product/usecases"
)

// OrderCase encapsulates the business logic related to orders.
type OrderCase struct {
	gateway     g.OrderGateway
	productCase *pc.ProductCase
}

// New creates a new instance of OrderCase with the provided OrderGateway.
func New(gateway g.OrderGateway, productCase *pc.ProductCase) *OrderCase {
	return &OrderCase{gateway, productCase}
}
