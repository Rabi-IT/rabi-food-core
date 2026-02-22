package order_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormOrderGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
