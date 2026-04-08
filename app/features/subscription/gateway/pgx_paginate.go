package subscription_gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxSubscriptionGatewayAdapter) Paginate(
	ctx context.Context,
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("subscriptions").PlaceholderFormat(sq.Dollar)

	if filter.TenantID != nil {
		base = base.Where(sq.Eq{"tenant_id": *filter.TenantID})
	}

	if filter.UserID != nil {
		base = base.Where(sq.Eq{"user_id": *filter.UserID})
	}

	if filter.Status != "" {
		base = base.Where(sq.Eq{"status": filter.Status})
	}

	countSQL, countArgs, err := base.Column("COUNT(*)").ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	dataSQL, dataArgs, err := base.
		Columns("id", "tenant_id", "user_id", "status", "total_cycles", "remaining_cycles",
			"payment_amount", "auto_renew", "created_at").
		Limit(uint64(paginate.PageSize)).
		Offset(uint64(paginate.Offset())).
		ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	data, count, err := database.Paginate[PaginateData](
		ctx, g.DB.Pool,
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
