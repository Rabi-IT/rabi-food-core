package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrForbidden        = newErr("forbidden", http.StatusForbidden)
	ErrNoValuesToUpdate = newErr("NO_VALUES_TO_UPDATE", http.StatusBadRequest)
	ErrConflict         = newErr("CONFLICT", http.StatusConflict)
)

type ApiError struct {
	Code     string `json:"code"`
	FullCode string `json:"fullCode,omitempty"`
	Status   int    `json:"status"`
	err      error  `json:"-"`
}

func newErr(code string, status int) *ApiError {
	return &ApiError{Code: code, Status: status}
}

func (e *ApiError) fullCode() string {
	if e.FullCode != "" {
		return e.FullCode
	}

	return e.Code
}

func (e *ApiError) Error() string {
	code := e.fullCode()

	if e.err != nil {
		return fmt.Sprintf("%s: %v", code, e.err)
	}

	return code
}

func (e *ApiError) Unwrap() error {
	return e.err
}

func (e *ApiError) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"code":"%s","status":%d}`, e.fullCode(), e.Status), nil
}
