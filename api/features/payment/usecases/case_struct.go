package payment_usecases

import (
	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
)

type PaymentCase struct {
	stripe           g.StripeGateway
	customer         g.StripeCustomerGateway
	subscriptionCase *subscription_usecases.SubscriptionCase
}

func New(
	stripe g.StripeGateway,
	customer g.StripeCustomerGateway,
	subscriptionCase *subscription_usecases.SubscriptionCase,
) *PaymentCase {
	return &PaymentCase{
		stripe:           stripe,
		customer:         customer,
		subscriptionCase: subscriptionCase,
	}
}
