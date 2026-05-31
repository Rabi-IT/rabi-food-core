package payment_gateway

import (
	"context"
	"fmt"

	auth_gateway "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	stripe "github.com/stripe/stripe-go/v82"
)

var offSession = "off_session"

func (g *StripePaymentGatewayAdapter) ChargeWithToken(ctx context.Context, input ChargeWithTokenInput) (ChargeOutput, error) {
	customerID, err := g.ensureCustomer(ctx, input.UserID, input.UserEmail, input.UserName)
	if err != nil {
		return ChargeOutput{}, err
	}

	params := &stripe.PaymentIntentCreateParams{
		Amount:            new(int64(input.AmountCents)),
		Currency:          &input.Currency,
		Customer:          &customerID,
		ConfirmationToken: &input.PaymentToken,
		Confirm:           new(true),
		SetupFutureUsage:  &offSession,
		AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
			Enabled:        new(true),
			AllowRedirects: stripe.String("never"),
		},
		Metadata: map[string]string{
			"subscription_id": input.SubscriptionID,
			"tenant_id":       input.TenantID,
		},
	}

	params.SetIdempotencyKey(input.IdempotencyKey)

	pi, err := g.sc.V1PaymentIntents.Create(ctx, params)
	if err != nil {
		return ChargeOutput{}, fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	if pi.Status == stripe.PaymentIntentStatusRequiresAction {
		return ChargeOutput{
			PaymentIntentID: pi.ID,
			Status:          string(pi.Status),
			ClientSecret:    pi.ClientSecret,
		}, nil
	}

	if pi.PaymentMethod != nil {
		pmID := pi.PaymentMethod.ID
		if _, err = g.auth.Patch(ctx, input.UserID, auth_gateway.PatchInput{
			StripePaymentMethodID: &pmID,
		}); err != nil {
			return ChargeOutput{}, fmt.Errorf("failed to save payment method: %w", err)
		}
	}

	return ChargeOutput{
		PaymentIntentID: pi.ID,
		Status:          string(pi.Status),
	}, nil
}

func (g *StripePaymentGatewayAdapter) ensureCustomer(ctx context.Context, userID, email, name string) (string, error) {
	user, err := g.auth.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to look up customer: %w", err)
	}

	if user != nil && user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	params := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}
	params.SetIdempotencyKey("customer-" + userID)

	stripeCustomer, err := g.sc.V1Customers.Create(ctx, params)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrStripeServiceFailure, err)
	}

	if _, err = g.auth.Patch(ctx, userID, auth_gateway.PatchInput{
		StripeCustomerID: &stripeCustomer.ID,
	}); err != nil {
		return "", fmt.Errorf("failed to save customer: %w", err)
	}

	return stripeCustomer.ID, nil
}
