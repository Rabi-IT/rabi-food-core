package payment_usecases

import (
	"context"
	"fmt"

	g "github.com/Rabi-IT/rabi-food-core/features/payment/gateway"
	subscription_usecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/google/uuid"
)

const paymentCurrency = "brl"

type CreateIntentInput struct {
	TenantID    string
	UserID      string
	UserEmail   string
	UserName    string
	Items       []subscription_usecases.SubscriptionItemInput
	TotalCycles uint
}

type CreateIntentOutput struct {
	ClientSecret string `json:"clientSecret"`
}

func (c *PaymentCase) CreateIntent(ctx context.Context, input CreateIntentInput) (CreateIntentOutput, error) {
	totalCents, err := c.subscriptionCase.CalculateTotalAmount(ctx, input.TenantID, input.Items, input.TotalCycles)
	if err != nil {
		return CreateIntentOutput{}, fmt.Errorf("failed to calculate amount: %w", err)
	}

	customerID, err := c.ensureStripeCustomer(ctx, input.UserID, input.UserEmail, input.UserName)
	if err != nil {
		return CreateIntentOutput{}, err
	}

	paymentAttemptID := uuid.Must(uuid.NewV7()).String()

	pi, err := c.stripe.CreatePaymentIntent(ctx, g.CreateIntentInput{
		CustomerID:     customerID,
		AmountCents:    int64(totalCents),
		Currency:       paymentCurrency,
		IdempotencyKey: "pi-" + paymentAttemptID,
	})
	if err != nil {
		return CreateIntentOutput{}, err
	}

	return CreateIntentOutput{ClientSecret: pi.ClientSecret}, nil
}

func (c *PaymentCase) ensureStripeCustomer(ctx context.Context, userID, email, name string) (string, error) {
	record, err := c.customer.GetByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to look up stripe customer: %w", err)
	}

	if record != nil {
		return record.StripeCustomerID, nil
	}

	stripeCustomerID, err := c.stripe.CreateCustomer(ctx, g.CreateStripeCustomerInput{
		UserID:         userID,
		Email:          email,
		Name:           name,
		IdempotencyKey: "customer-" + userID,
	})
	if err != nil {
		return "", err
	}

	if err = c.customer.Upsert(ctx, g.UpsertCustomerInput{
		UserID:           userID,
		StripeCustomerID: stripeCustomerID,
		HasPaymentMethod: false,
	}); err != nil {
		return "", fmt.Errorf("failed to save stripe customer: %w", err)
	}

	return stripeCustomerID, nil
}
