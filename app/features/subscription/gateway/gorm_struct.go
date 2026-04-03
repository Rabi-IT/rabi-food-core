package subscription_gateway

import "rabi-food-core/libs/database"

type GormSubscriptionGatewayAdapter struct {
	DB *database.GormAdapter
}
