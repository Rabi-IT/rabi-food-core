package payment_gateway

import (
	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	stripe "github.com/stripe/stripe-go/v82"
)

type StripePaymentGatewayAdapter struct {
	sc *stripe.Client
	db *database.PgxAdapter
}

func NewStripe(db *database.PgxAdapter) PaymentGateway {
	stripe.EnableTelemetry = false

	return &StripePaymentGatewayAdapter{
		sc: stripe.NewClient(config.StripeSecretKey),
		db: db,
	}
}
