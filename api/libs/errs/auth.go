package errs

import "net/http"

var (
	ErrInvalidCredentials  = newErr("INVALID_CREDENTIALS", http.StatusUnauthorized)
	ErrInvalidRefreshToken = newErr("INVALID_REFRESH_TOKEN", http.StatusUnauthorized)
	ErrInvalidOTP          = newErr("INVALID_OTP", http.StatusUnauthorized)
	ErrInvalidToken        = newErr("INVALID_TOKEN", http.StatusUnauthorized)
	ErrUserAlreadyExists   = newErr("USER_ALREADY_EXISTS", http.StatusConflict)
	ErrAuthServiceFailure  = newErr("AUTH_SERVICE_FAILURE", http.StatusServiceUnavailable)
	ErrMissingTenantID     = newErr("MISSING_TENANT_ID", http.StatusBadRequest)
)
