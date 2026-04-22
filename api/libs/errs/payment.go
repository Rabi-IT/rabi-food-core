package errs

import "net/http"

var (
	// Payment already confirmed with a different external_payment_id.
	ErrPaymentExternalIDConflict = newErr("PAYMENT_EXTERNAL_ID_CONFLICT", http.StatusConflict)
	ErrUnknownBusinessReason     = newErr("UNKNOWN_BUSINESS_REASON", http.StatusInternalServerError)
)
