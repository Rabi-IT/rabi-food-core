package subscription_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormSubscriptionGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
