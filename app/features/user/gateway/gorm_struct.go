package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormUserGatewayAdapter struct {
	DB *database.GormAdapter
}
