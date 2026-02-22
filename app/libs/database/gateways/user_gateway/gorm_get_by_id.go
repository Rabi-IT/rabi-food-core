package user_gateway

import (
	"errors"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

var (
	ErrMustFilterByIDOrTenantID = errors.New("must filter by either ID or TenantID")
)

func (g *GormUserGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	if filter.ID == "" && filter.TenantID == "" {
		return nil, ErrMustFilterByIDOrTenantID
	}

	output := &models.User{}
	query := g.DB.Conn.Limit(1)

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Find(output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:         output.ID,
		State:      output.State,
		ZIP:        output.ZIP,
		Phone:      output.Phone,
		City:       output.City,
		Photo:      output.Photo,
		TaxID:      output.TaxID,
		SocialID:   output.SocialID,
		Street:     output.Street,
		Complement: output.Complement,
		Name:       output.Name,
		Email:      output.Email,
		TenantID:   output.TenantID,
	}

	return &adapted, nil
}
