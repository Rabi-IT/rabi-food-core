package usecases

import (
	pc "rabi-food-core/features/product/usecases"
	g "rabi-food-core/features/subscription/gateway"
)

// SubscriptionCase encapsulates the business logic related to subscriptions.
type SubscriptionCase struct {
	gateway     g.SubscriptionGateway
	productCase *pc.ProductCase
}

// New creates a new instance of SubscriptionCase.
func New(gateway g.SubscriptionGateway, productCase *pc.ProductCase) *SubscriptionCase {
	return &SubscriptionCase{gateway, productCase}
}
