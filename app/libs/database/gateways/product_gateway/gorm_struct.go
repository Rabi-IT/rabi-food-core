package product_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormProductGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
