package payment_gateway

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
)

type userPaymentMeta struct {
	StripeCustomerID      string
	StripePaymentMethodID string
}

func (g *StripePaymentGatewayAdapter) getUserMeta(ctx context.Context, userID string) (userPaymentMeta, error) {
	var customerID, pmID string

	err := g.db.Pool.QueryRow(ctx,
		`SELECT COALESCE(raw_app_meta_data->>'stripe_customer_id', ''),
		        COALESCE(raw_app_meta_data->>'stripe_payment_method_id', '')
		 FROM auth.users WHERE id = $1`,
		userID,
	).Scan(&customerID, &pmID)

	if errors.Is(err, pgx.ErrNoRows) {
		return userPaymentMeta{}, nil
	}

	return userPaymentMeta{StripeCustomerID: customerID, StripePaymentMethodID: pmID}, err
}

func (g *StripePaymentGatewayAdapter) saveUserMeta(ctx context.Context, userID string, patch map[string]string) error {
	data, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = g.db.Pool.Exec(context.WithoutCancel(ctx),
		`UPDATE auth.users SET raw_app_meta_data = raw_app_meta_data || $1::jsonb WHERE id = $2`,
		string(data), userID,
	)

	return err
}
