package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormCategoryGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.NewString()

	result := g.DB.Conn.Create(&models.Category{
		ID:          id,
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
	})

	return id, result.Error
}
