package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxOrderGatewayAdapter struct {
	DB *database.PgxAdapter
}
