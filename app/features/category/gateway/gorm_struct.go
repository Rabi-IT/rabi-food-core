package gateway

import "rabi-food-core/libs/database"

type GormCategoryGatewayAdapter struct {
	DB *database.GormAdapter
}
