package errs

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrOrderNotFound                = newErr("ORDER_NOT_FOUND", http.StatusNotFound)
	ErrProductNotFound              = newErr("PRODUCT_NOT_FOUND", http.StatusNotFound)
	ErrStatusNotModifiable          = newErr("STATUS_NOT_MODIFIABLE", http.StatusBadRequest)
	ErrInvalidStatusTransition      = newErr("INVALID_STATUS_TRANSITION", http.StatusBadRequest)
	ErrOrderStateVerificationFailed = newErr("ORDER_STATE_VERIFICATION_FAILED", http.StatusInternalServerError)
)

func ProductNotFound(productIDs ...string) *ApiError {
	if len(productIDs) == 0 {
		return ErrProductNotFound
	}

	e := *ErrProductNotFound
	e.Code = fmt.Sprintf("%s__%s", ErrProductNotFound.Code, strings.Join(productIDs, "_"))

	// To allow errors.Is checks
	e.err = ErrProductNotFound

	return &e
}

func StatusNotModifiable(status fmt.Stringer) *ApiError {
	e := *ErrStatusNotModifiable
	e.Code = ErrStatusNotModifiable.Code
	e.FullCode = fmt.Sprintf("%s__%s", e.Code, status)

	// To allow errors.Is checks
	e.err = ErrStatusNotModifiable

	return &e
}

func InvalidTranstion(from, to fmt.Stringer) *ApiError {
	e := *ErrInvalidStatusTransition
	e.Code = ErrInvalidStatusTransition.Code
	e.FullCode = fmt.Sprintf("%s__FROM_%s_TO_%s", e.Code, from, to)

	// To allow errors.Is checks
	e.err = ErrInvalidStatusTransition

	return &e
}
