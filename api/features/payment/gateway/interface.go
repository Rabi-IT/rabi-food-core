package payment_gateway

import "context"

type PaymentGateway interface {
	ChargeWithToken(ctx context.Context, input ChargeWithTokenInput) (ChargeOutput, error)
	ChargeOffSession(ctx context.Context, input ChargeOffSessionInput) (ChargeOutput, error)
	GetPaymentProfile(ctx context.Context, userID string) (*PaymentProfileOutput, error)
	Refund(ctx context.Context, paymentIntentID string) error
}

type ChargeWithTokenInput struct {
	UserID         string
	UserEmail      string
	UserName       string
	PaymentToken   string
	AmountCents    int64
	Currency       string
	Metadata       map[string]string
	IdempotencyKey string
}

type ChargeOffSessionInput struct {
	UserID         string
	AmountCents    int64
	Currency       string
	Metadata       map[string]string
	IdempotencyKey string
}

type ChargeOutput struct {
	PaymentIntentID string
	Status          string
	ClientSecret    string // non-empty when requires_action
}

type PaymentProfileOutput struct {
	HasPaymentMethod bool
	Brand            string
	Last4            string
	ExpMonth         int64
	ExpYear          int64
}
