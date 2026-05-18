package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) Patch(ctx context.Context, id string, input PatchInput) (bool, error) {
	qb := sq.Update("auth.users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	if input.Email != nil {
		qb = qb.Set("email", *input.Email)
	}

	appMeta := map[string]any{}

	if input.StripeCustomerID != nil {
		appMeta["stripe_customer_id"] = *input.StripeCustomerID
	}
	if input.StripePaymentMethodID != nil {
		appMeta["stripe_payment_method_id"] = *input.StripePaymentMethodID
	}

	if len(appMeta) > 0 {
		data, err := json.Marshal(appMeta)
		if err != nil {
			return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
		}
		qb = qb.Set("raw_app_meta_data", sq.Expr("raw_app_meta_data || ?::jsonb", string(data)))
	}

	userMeta := map[string]any{}

	if input.Name != nil {
		userMeta["name"] = *input.Name
	}
	if input.Phone != nil {
		userMeta["phone"] = *input.Phone
	}
	if input.TaxID != nil {
		userMeta["tax_id"] = *input.TaxID
	}
	if input.SocialID != nil {
		userMeta["social_id"] = *input.SocialID
	}
	if input.City != nil {
		userMeta["city"] = *input.City
	}
	if input.State != nil {
		userMeta["state"] = *input.State
	}
	if input.ZIP != nil {
		userMeta["zip"] = *input.ZIP
	}
	if input.Street != nil {
		userMeta["street"] = *input.Street
	}
	if input.Complement != nil {
		userMeta["complement"] = *input.Complement
	}
	if input.Neighborhood != nil {
		userMeta["neighborhood"] = *input.Neighborhood
	}

	if len(userMeta) > 0 {
		data, err := json.Marshal(userMeta)
		if err != nil {
			return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
		}
		qb = qb.Set("raw_user_meta_data", sq.Expr("raw_user_meta_data || ?::jsonb", string(data)))
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	tag, err := g.db.Pool.Exec(context.WithoutCancel(ctx), sql, args...)
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	return tag.RowsAffected() > 0, nil
}
