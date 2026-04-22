package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxTenantGatewayAdapter struct {
	DB *database.PgxAdapter
}
