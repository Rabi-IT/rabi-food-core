package subscription_case

import (
	g "rabi-food-core/libs/database/gateways/subscription_gateway"
	pc "rabi-food-core/usecases/product_case"
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
