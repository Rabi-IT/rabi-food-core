package errs

import "net/http"

var (
	ErrInvalidCredentials  = newErr("INVALID_CREDENTIALS", http.StatusUnauthorized)
	ErrInvalidRefreshToken = newErr("INVALID_REFRESH_TOKEN", http.StatusUnauthorized)
	ErrAuthServiceFailure  = newErr("AUTH_SERVICE_FAILURE", http.StatusServiceUnavailable)
)
