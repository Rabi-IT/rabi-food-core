package payment_usecases

import (
	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
)

type PaymentCase struct {
	payment          g.PaymentGateway
	subscriptionCase *subscription_usecases.SubscriptionCase
}

func New(
	payment g.PaymentGateway,
	subscriptionCase *subscription_usecases.SubscriptionCase,
) *PaymentCase {
	return &PaymentCase{
		payment:          payment,
		subscriptionCase: subscriptionCase,
	}
}
