package payment_gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxStripeCustomerAdapter) GetByUserID(ctx context.Context, userID string) (*CustomerOutput, error) {
	sql, args, err := sq.
		Select("user_id", "stripe_customer_id", "COALESCE(stripe_payment_method_id, '')", "has_payment_method").
		From("payment.stripe_customers").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var out CustomerOutput
	err = g.DB.Pool.QueryRow(ctx, sql, args...).
		Scan(&out.UserID, &out.StripeCustomerID, &out.StripePaymentMethodID, &out.HasPaymentMethod)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &out, err
}
