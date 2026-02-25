package user_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormUserGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	result := g.DB.Conn.Create(&models.User{
		ID:           id,
		TenantID:     input.TenantID,
		TaxID:        input.TaxID,
		SocialID:     input.SocialID,
		Street:       input.Street,
		Complement:   input.Complement,
		Name:         input.Name,
		Email:        input.Email,
		Photo:        input.Photo,
		ZIP:          input.ZIP,
		Phone:        input.Phone,
		City:         input.City,
		State:        input.State,
		Neighborhood: input.Neighborhood,
		Role:         input.Role,
	})

	return id, result.Error
}
