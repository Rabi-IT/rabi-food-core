package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

var (
	ErrForbidden        = newErr("forbidden", http.StatusForbidden)
	ErrNoValuesToUpdate = newErr("NO_VALUES_TO_UPDATE", http.StatusBadRequest)
	ErrConflict         = newErr("CONFLICT", http.StatusConflict)
)

type AppError struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	err    error  `json:"-"`
}

func newErr(code string, status int) *AppError {
	return &AppError{Code: code, Status: status}
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.err)
	}

	return e.Code
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"code":"%s","status":%d}`, e.Code, e.Status), nil
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505"
	}

	return false
}
