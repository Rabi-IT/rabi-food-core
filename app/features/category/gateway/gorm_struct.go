package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormCategoryGatewayAdapter struct {
	DB *database.GormAdapter
}
