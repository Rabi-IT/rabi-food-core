package order_gateway

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormOrderGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	query := g.DB.Conn.Model(&models.Order{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", filter.CreatedAtFrom)
	}

	if filter.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", filter.CreatedAtTo)
	}

	if filter.FulfillmentStatus != "" {
		query = query.Where("fulfillment_status = ?", filter.FulfillmentStatus)
	}

	if filter.DeliveryStatus != "" {
		query = query.Where("delivery_status = ?", filter.DeliveryStatus)
	}

	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", filter.TenantID)
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
