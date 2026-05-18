package errs

import "net/http"

var (
	ErrPaymentExternalIDConflict  = newErr("PAYMENT_EXTERNAL_ID_CONFLICT", http.StatusConflict)
	ErrUnknownBusinessReason      = newErr("UNKNOWN_BUSINESS_REASON", http.StatusInternalServerError)
	ErrPaymentIntentNotSucceeded  = newErr("PAYMENT_INTENT_NOT_SUCCEEDED", http.StatusPaymentRequired)
	ErrPaymentAlreadyProcessed    = newErr("PAYMENT_ALREADY_PROCESSED", http.StatusConflict)
	ErrNoPaymentMethod            = newErr("NO_PAYMENT_METHOD", http.StatusPaymentRequired)
	ErrStripeServiceFailure       = newErr("STRIPE_SERVICE_FAILURE", http.StatusBadGateway)
	ErrPaymentRequiresAction      = newErr("PAYMENT_REQUIRES_ACTION", http.StatusPaymentRequired)
)
