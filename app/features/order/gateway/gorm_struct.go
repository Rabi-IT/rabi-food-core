package gateway

import "rabi-food-core/libs/database"

type GormOrderGatewayAdapter struct {
	DB *database.GormAdapter
}
