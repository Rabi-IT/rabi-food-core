package app_plan_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormAppPlanGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
