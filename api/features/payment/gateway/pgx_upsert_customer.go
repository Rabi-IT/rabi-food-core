package payment_gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxStripeCustomerAdapter) Upsert(ctx context.Context, input UpsertCustomerInput) error {
	sql, args, err := sq.
		Insert("payment.stripe_customers").
		Columns("user_id", "stripe_customer_id", "stripe_payment_method_id", "has_payment_method").
		Values(input.UserID, input.StripeCustomerID, nullableString(input.StripePaymentMethodID), input.HasPaymentMethod).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET stripe_customer_id = EXCLUDED.stripe_customer_id, stripe_payment_method_id = EXCLUDED.stripe_payment_method_id, has_payment_method = EXCLUDED.has_payment_method, updated_at = NOW()").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return err
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}

	return s
}
