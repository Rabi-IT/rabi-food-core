package tenant_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormTenantGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	result := g.DB.Conn.Create(&models.Tenant{
		ID:   id,
		Name: input.Name,
	})

	return id, result.Error
}
