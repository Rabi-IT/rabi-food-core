package subscription_gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormSubscriptionGatewayAdapter struct {
	DB *database.GormAdapter
}
