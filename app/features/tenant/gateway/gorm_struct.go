package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type GormTenantGatewayAdapter struct {
	DB *database.GormAdapter
}
