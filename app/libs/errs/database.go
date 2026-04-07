package errs

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
