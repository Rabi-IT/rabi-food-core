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

type AppError struct {
	Code     string `json:"code"`
	FullCode string `json:"fullCode,omitempty"`
	Status   int    `json:"status"`
	err      error  `json:"-"`
}

func newErr(code string, status int) *AppError {
	return &AppError{Code: code, Status: status}
}

func (e *AppError) fullCode() string {
	if e.FullCode != "" {
		return e.FullCode
	}

	return e.Code
}

func (e *AppError) Error() string {
	code := e.fullCode()

	if e.err != nil {
		return fmt.Sprintf("%s: %v", code, e.err)
	}

	return code
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"code":"%s","status":%d}`, e.fullCode(), e.Status), nil
}
