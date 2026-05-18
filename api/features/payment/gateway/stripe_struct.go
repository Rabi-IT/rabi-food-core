package payment_gateway

import (
	"github.com/Rabi-IT/rabi-food-core/config"
	auth_gateway "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	stripe "github.com/stripe/stripe-go/v82"
)

type StripePaymentGatewayAdapter struct {
	sc   *stripe.Client
	auth auth_gateway.AuthGateway
}

func NewStripe(auth auth_gateway.AuthGateway) PaymentGateway {
	stripe.EnableTelemetry = false

	return &StripePaymentGatewayAdapter{
		sc:   stripe.NewClient(config.StripeSecretKey),
		auth: auth,
	}
}
