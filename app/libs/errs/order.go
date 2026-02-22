package errs

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrProductNotFound         = newErr("PRODUCT_NOT_FOUND", http.StatusNotFound)
	ErrStatusNotModifiable     = newErr("STATUS_NOT_MODIFIABLE", http.StatusBadRequest)
	ErrInvalidStatusTransition = newErr("INVALID_STATUS_TRANSITION", http.StatusBadRequest)
)

func ProductNotFound(productIDs ...string) *AppError {
	if len(productIDs) == 0 {
		return ErrProductNotFound
	}

	e := *ErrProductNotFound
	e.Code = fmt.Sprintf("%s__%s", ErrProductNotFound.Code, strings.Join(productIDs, "_"))

	// To allow errors.Is checks
	e.err = ErrProductNotFound

	return &e
}

func StatusNotModifiable(status fmt.Stringer) *AppError {
	e := *ErrStatusNotModifiable
	e.Code = fmt.Sprintf("%s__%s", ErrStatusNotModifiable.Code, status)

	// To allow errors.Is checks
	e.err = ErrStatusNotModifiable

	return &e
}

func InvalidTranstion(from, to fmt.Stringer) *AppError {
	e := *ErrInvalidStatusTransition
	e.Code = fmt.Sprintf("%s__FROM_%s_TO_%s", ErrInvalidStatusTransition.Code, from, to)

	// To allow errors.Is checks
	e.err = ErrInvalidStatusTransition

	return &e
}
