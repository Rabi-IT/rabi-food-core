package user_gateway

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormUserGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	query := g.DB.Conn.Model(&models.User{})

	if filter.Name != nil {
		query = query.Where("name = ?", filter.Name)
	}

	if filter.City != nil {
		query = query.Where("city = ?", filter.City)
	}

	if filter.State != nil {
		query = query.Where("state = ?", filter.State)
	}

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}

	count := int64(0)
	data := []PaginateData{}
	err := gorm_adapter.Paginate(query, &count, &data, paginate)
	if err != nil {
		return PaginateOutput{}, err
	}

	output := PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}

	return output, nil
}
