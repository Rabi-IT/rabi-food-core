package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNoValuesToUpdate = newErr("NO_VALUES_TO_UPDATE", http.StatusBadRequest)
	ErrConflict         = newErr("CONFLICT", http.StatusConflict)
)

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
