package usecases

import (
	payment_gateway "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	pc "github.com/Rabi-IT/rabi-food-core/features/product/usecases"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
)

type SubscriptionCase struct {
	gateway        g.SubscriptionGateway
	productCase    *pc.ProductCase
	paymentGateway payment_gateway.PaymentGateway
}

func New(gateway g.SubscriptionGateway, productCase *pc.ProductCase, paymentGateway payment_gateway.PaymentGateway) *SubscriptionCase {
	return &SubscriptionCase{gateway, productCase, paymentGateway}
}
