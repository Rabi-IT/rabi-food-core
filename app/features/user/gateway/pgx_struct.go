package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxUserGatewayAdapter struct {
	DB *database.PgxAdapter
}
