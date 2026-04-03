package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormProductGatewayAdapter struct {
	DB *database.GormAdapter
}
