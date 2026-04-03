package usecases

import g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"

// ProductCase encapsulates the business logic related to products.
type ProductCase struct {
	gateway g.ProductGateway
}

// New creates a new instance of ProductCase with the provided ProductGateway.
func New(gateway g.ProductGateway) *ProductCase {
	return &ProductCase{gateway}
}
