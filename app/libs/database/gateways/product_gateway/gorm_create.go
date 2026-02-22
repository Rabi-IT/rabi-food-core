package product_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormProductGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.NewString()

	result := g.DB.Conn.Create(&models.Product{
		ID:          id,
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
		Photo:       input.Photo,
		CategoryID:  input.CategoryID,
		Unit:        input.Unit,
		Price:       input.Price,
		IsActive:    input.IsActive,
	})

	return id, result.Error
}
