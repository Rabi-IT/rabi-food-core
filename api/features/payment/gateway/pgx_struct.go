package payment_gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxStripeCustomerAdapter struct {
	DB *database.PgxAdapter
}
