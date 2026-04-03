package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormOrderGatewayAdapter struct {
	DB *database.GormAdapter
}
