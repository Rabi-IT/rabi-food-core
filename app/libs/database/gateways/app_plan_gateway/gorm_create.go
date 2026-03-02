package app_plan_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormAppPlanGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	result := g.DB.Conn.Create(&models.AppPlan{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		IsActive:    input.IsActive,
	})

	return id, result.Error
}
