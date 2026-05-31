package payment_gateway

import "context"

// PaymentGateway handles all communication with the payment provider.
//
// ChargeWithToken is used for the first payment — the user provides a confirmation
// token generated on the frontend. On success, the card is saved for future charges.
//
// ChargeOffSession is used for recurring charges against a previously saved card,
// without the user being present (e.g. subscription renewal).
type PaymentGateway interface {
	ChargeWithToken(ctx context.Context, input ChargeWithTokenInput) (ChargeOutput, error)
	ChargeOffSession(ctx context.Context, input ChargeOffSessionInput) (ChargeOutput, error)
	GetProfile(ctx context.Context, userID string) (*GetProfileOutput, error)
	Refund(ctx context.Context, paymentIntentID string) error
}

type ChargeWithTokenInput struct {
	UserID         string
	UserEmail      string
	UserName       string
	PaymentToken   string
	SubscriptionID string
	TenantID       string
	AmountCents    uint
	Currency       string
	IdempotencyKey string
}

type ChargeOffSessionInput struct {
	UserID         string
	SubscriptionID string
	TenantID       string
	AmountCents    uint
	Currency       string
	IdempotencyKey string
}

type ChargeOutput struct {
	PaymentIntentID string
	Status          string
	ClientSecret    string // non-empty when requires_action
}

type GetProfileOutput struct {
	HasPaymentMethod bool
	Brand            string
	Last4            string
	ExpMonth         int64
	ExpYear          int64
}
