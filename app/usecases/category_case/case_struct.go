package category_case

import g "rabi-food-core/libs/database/gateways/category_gateway"

// CategoryCase encapsulates the business logic related to categorys.
type CategoryCase struct {
	gateway g.CategoryGateway
}

// New creates a new instance of CategoryCase with the provided CategoryGateway.
func New(gateway g.CategoryGateway) *CategoryCase {
	return &CategoryCase{gateway}
}
