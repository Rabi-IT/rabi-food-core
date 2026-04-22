package subscription_gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxSubscriptionGatewayAdapter struct {
	DB *database.PgxAdapter
}
