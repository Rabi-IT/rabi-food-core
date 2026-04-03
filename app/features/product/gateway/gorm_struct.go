package gateway

import "rabi-food-core/libs/database"

type GormProductGatewayAdapter struct {
	DB *database.GormAdapter
}
