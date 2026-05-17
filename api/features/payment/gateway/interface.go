package payment_gateway

import "context"

// StripeGateway handles all communication with the Stripe API.
type StripeGateway interface {
	CreateCustomer(ctx context.Context, input CreateStripeCustomerInput) (string, error)
	CreatePaymentIntent(ctx context.Context, input CreateIntentInput) (CreateIntentOutput, error)
	GetPaymentIntent(ctx context.Context, id string) (*PaymentIntentOutput, error)
	GetPaymentMethod(ctx context.Context, id string) (*PaymentMethodOutput, error)
	ConfirmOffSession(ctx context.Context, input ConfirmOffSessionInput) (ConfirmOffSessionOutput, error)
	Refund(ctx context.Context, paymentIntentID string) error
}

// StripeCustomerGateway handles persistence of Stripe customer data.
type StripeCustomerGateway interface {
	Upsert(ctx context.Context, input UpsertCustomerInput) error
	GetByUserID(ctx context.Context, userID string) (*CustomerOutput, error)
}

type CreateStripeCustomerInput struct {
	UserID         string
	Email          string
	Name           string
	IdempotencyKey string
}

type CreateIntentInput struct {
	CustomerID     string
	AmountCents    int64
	Currency       string
	IdempotencyKey string
	Metadata       map[string]string
}

type CreateIntentOutput struct {
	PaymentIntentID string
	ClientSecret    string
}

type PaymentIntentOutput struct {
	ID              string
	Status          string
	PaymentMethodID string
}

type ConfirmOffSessionInput struct {
	CustomerID      string
	PaymentMethodID string
	AmountCents     int64
	Currency        string
	IdempotencyKey  string
	Metadata        map[string]string
}

type ConfirmOffSessionOutput struct {
	PaymentIntentID string
	Status          string
	ClientSecret    string // non-empty when requires_action
}

type UpsertCustomerInput struct {
	UserID                string
	StripeCustomerID      string
	StripePaymentMethodID string // empty if no card saved yet
	HasPaymentMethod      bool
}

type CustomerOutput struct {
	UserID                string
	StripeCustomerID      string
	StripePaymentMethodID string
	HasPaymentMethod      bool
}

type PaymentMethodOutput struct {
	Brand    string
	Last4    string
	ExpMonth int64
	ExpYear  int64
}
