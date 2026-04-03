package gateway

import "rabi-food-core/libs/database"

type GormUserGatewayAdapter struct {
	DB *database.GormAdapter
}
