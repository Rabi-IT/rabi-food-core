package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxProductGatewayAdapter struct {
	DB *database.PgxAdapter
}
