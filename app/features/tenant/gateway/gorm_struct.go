package gateway

import "rabi-food-core/libs/database"

type GormTenantGatewayAdapter struct {
	DB *database.GormAdapter
}
