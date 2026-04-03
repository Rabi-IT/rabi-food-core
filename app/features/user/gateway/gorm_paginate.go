package gateway

import (
	"github.com/Rabi-IT/rabi-food-core/features/user/model"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *GormUserGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	query := g.DB.Conn.Model(&model.User{})

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
	err := database.Paginate(query, &count, &data, paginate)
	if err != nil {
		return PaginateOutput{}, err
	}

	output := PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}

	return output, nil
}
