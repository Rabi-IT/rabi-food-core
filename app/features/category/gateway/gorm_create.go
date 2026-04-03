package gateway

import (
	"github.com/Rabi-IT/rabi-food-core/features/category/model"

	"github.com/google/uuid"
)

func (g *GormCategoryGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	result := g.DB.Conn.Create(&model.Category{
		ID:          id,
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
	})

	return id, result.Error
}
