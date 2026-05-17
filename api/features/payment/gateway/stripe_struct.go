package payment_gateway

import (
	"github.com/Rabi-IT/rabi-food-core/config"
	stripe "github.com/stripe/stripe-go/v82"
)

type StripeGatewayAdapter struct {
	sc *stripe.Client
}

func NewStripe() StripeGateway {
	stripe.EnableTelemetry = false

	return &StripeGatewayAdapter{
		sc: stripe.NewClient(config.StripeSecretKey),
	}
}
