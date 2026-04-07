package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxOrderGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("orders").
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != nil {
		base = base.Where(sq.Eq{"user_id": *filter.UserID})
	}

	if filter.CreatedAtFrom != nil {
		base = base.Where(sq.GtOrEq{"created_at": *filter.CreatedAtFrom})
	}

	if filter.CreatedAtTo != nil {
		base = base.Where(sq.LtOrEq{"created_at": *filter.CreatedAtTo})
	}

	if filter.FulfillmentStatus != "" {
		base = base.Where(sq.Eq{"fulfillment_status": filter.FulfillmentStatus})
	}

	if filter.DeliveryStatus != "" {
		base = base.Where(sq.Eq{"delivery_status": filter.DeliveryStatus})
	}

	if filter.PaymentStatus != "" {
		base = base.Where(sq.Eq{"payment_status": filter.PaymentStatus})
	}

	if filter.TenantID != nil {
		base = base.Where(sq.Eq{"tenant_id": *filter.TenantID})
	}

	countSQL, countArgs, err := base.Column("COUNT(*)").ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	dataSQL, dataArgs, err := base.
		Columns("id", "tenant_id", "code", "payment_status", "fulfillment_status",
			"delivery_status", "notes", "total_price", "created_at").
		Limit(uint64(paginate.PageSize)).
		Offset(uint64(paginate.Offset())).
		ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	data, count, err := database.Paginate[PaginateData](
		context.Background(), g.DB.Pool,
		countSQL, countArgs,
		dataSQL, dataArgs,
	)
	if err != nil {
		return PaginateOutput{}, err
	}

	return PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}, nil
}
