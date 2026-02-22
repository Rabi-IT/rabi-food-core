package order_case

import (
	g "rabi-food-core/libs/database/gateways/order_gateway"
	pc "rabi-food-core/usecases/product_case"
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
