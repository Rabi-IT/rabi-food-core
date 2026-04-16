package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxTenantGatewayAdapter) GetMember(ctx context.Context, filter GetMemberFilter) (*GetMemberOutput, error) {
	q := sq.
		Select("user_id", "tenant_id", "role").
		From("iam.tenant_members").
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != nil {
		q = q.Where(sq.Eq{"user_id": *filter.UserID})
	}

	if filter.TenantID != nil {
		q = q.Where(sq.Eq{"tenant_id": *filter.TenantID})
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var out GetMemberOutput
	err = g.DB.Pool.QueryRow(ctx, sql, args...).Scan(&out.UserID, &out.TenantID, &out.Role)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &out, err
}
