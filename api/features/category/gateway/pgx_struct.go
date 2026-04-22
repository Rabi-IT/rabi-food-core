package gateway

import "github.com/Rabi-IT/rabi-food-core/libs/database"

type PgxCategoryGatewayAdapter struct {
	DB *database.PgxAdapter
}
