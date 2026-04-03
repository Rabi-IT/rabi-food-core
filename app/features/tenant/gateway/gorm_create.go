package gateway

import (
	"github.com/Rabi-IT/rabi-food-core/features/tenant/model"

	"github.com/google/uuid"
)

func (g *GormTenantGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	result := g.DB.Conn.Create(&model.Tenant{
		ID:   id,
		Name: input.Name,
	})

	return id, result.Error
}
