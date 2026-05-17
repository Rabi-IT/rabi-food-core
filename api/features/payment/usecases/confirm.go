package payment_usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/google/uuid"
)

type ConfirmInput struct {
	TenantID        string
	UserID          string
	UserEmail       string
	UserName        string
	PaymentIntentID string // empty → charge saved card
	Items           []subscription_usecases.SubscriptionItemInput
	DeliveryDays    []subscription_gateway.DeliveryDay
	TotalCycles     uint
	Notes           string
}

type ConfirmOutput struct {
	SubscriptionIDs []string `json:"subscriptionIds"`
	ClientSecret    string   `json:"clientSecret,omitempty"`
}

func (c *PaymentCase) Confirm(ctx context.Context, input ConfirmInput) (ConfirmOutput, error) {
	if input.PaymentIntentID != "" {
		return c.confirmWithNewCard(ctx, input)
	}

	return c.confirmWithSavedCard(ctx, input)
}

func (c *PaymentCase) confirmWithNewCard(ctx context.Context, input ConfirmInput) (ConfirmOutput, error) {
	pi, err := c.stripe.GetPaymentIntent(ctx, input.PaymentIntentID)
	if err != nil {
		return ConfirmOutput{}, err
	}

	if pi.Status != "succeeded" {
		return ConfirmOutput{}, errs.ErrPaymentIntentNotSucceeded
	}

	customerID, err := c.ensureStripeCustomer(ctx, input.UserID, input.UserEmail, input.UserName)
	if err != nil {
		return ConfirmOutput{}, err
	}

	if err = c.customer.Upsert(ctx, g.UpsertCustomerInput{
		UserID:                input.UserID,
		StripeCustomerID:      customerID,
		StripePaymentMethodID: pi.PaymentMethodID,
		HasPaymentMethod:      true,
	}); err != nil {
		return ConfirmOutput{}, fmt.Errorf("failed to update payment method: %w", err)
	}

	ids, err := c.createSubscription(ctx, input)
	if err != nil {
		return ConfirmOutput{}, err
	}

	return ConfirmOutput{SubscriptionIDs: ids}, nil
}

func (c *PaymentCase) confirmWithSavedCard(ctx context.Context, input ConfirmInput) (ConfirmOutput, error) {
	record, err := c.customer.GetByUserID(ctx, input.UserID)
	if err != nil {
		return ConfirmOutput{}, fmt.Errorf("failed to look up customer: %w", err)
	}

	if record == nil || !record.HasPaymentMethod || record.StripePaymentMethodID == "" {
		return ConfirmOutput{}, errs.ErrNoPaymentMethod
	}

	totalCents, err := c.subscriptionCase.CalculateTotalAmount(ctx, input.TenantID, input.Items, input.TotalCycles)
	if err != nil {
		return ConfirmOutput{}, fmt.Errorf("failed to calculate amount: %w", err)
	}

	result, err := c.stripe.ConfirmOffSession(ctx, g.ConfirmOffSessionInput{
		CustomerID:      record.StripeCustomerID,
		PaymentMethodID: record.StripePaymentMethodID,
		AmountCents:     int64(totalCents),
		Currency:        paymentCurrency,
		IdempotencyKey:  "confirm-" + uuid.Must(uuid.NewV7()).String(),
	})
	if err != nil {
		return ConfirmOutput{}, err
	}

	if result.Status == "requires_action" {
		return ConfirmOutput{ClientSecret: result.ClientSecret}, errs.ErrPaymentRequiresAction
	}

	if result.Status != "succeeded" {
		return ConfirmOutput{}, errs.ErrPaymentIntentNotSucceeded
	}

	ids, err := c.createSubscription(ctx, input)
	if err != nil {
		return ConfirmOutput{}, err
	}

	return ConfirmOutput{SubscriptionIDs: ids}, nil
}

func (c *PaymentCase) createSubscription(ctx context.Context, input ConfirmInput) ([]string, error) {
	ids, err := c.subscriptionCase.Create(ctx, subscription_usecases.CreateInput{
		TenantID:     input.TenantID,
		UserID:       input.UserID,
		Items:        input.Items,
		DeliveryDays: input.DeliveryDays,
		TotalCycles:  input.TotalCycles,
		Notes:        input.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	sysCtx := app_context.WithSession(ctx, &app_context.UserSession{Role: auth.System})
	paidAt := time.Now().UTC()

	for _, id := range ids {
		_, err = c.subscriptionCase.ConfirmPayment(sysCtx, subscription_usecases.ConfirmPaymentInput{
			SubscriptionID: id,
			Provider:       "stripe",
			PaidAt:         paidAt,
		})
		if err != nil {
			return ids, fmt.Errorf("failed to confirm subscription payment: %w", err)
		}
	}

	return ids, nil
}
