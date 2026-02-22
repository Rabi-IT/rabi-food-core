package category_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormCategoryGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
