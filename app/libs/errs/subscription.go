package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrDuplicateProduct                = newErr("DUPLICATE_PRODUCT", http.StatusBadRequest)
	ErrDuplicateDeliveryWeekday        = newErr("DUPLICATE_DELIVERY_WEEKDAY", http.StatusBadRequest)
	ErrTenantSubscriptionNotConfigured = newErr("TENANT_SUBSCRIPTION_NOT_CONFIGURED", http.StatusInternalServerError)
)

func DuplicateProduct(productID string) *AppError {
	e := *ErrDuplicateProduct
	e.Code = ErrDuplicateProduct.Code
	e.FullCode = fmt.Sprintf("%s__%s", e.Code, productID)

	// allow errors.Is
	e.err = ErrDuplicateProduct

	return &e
}

func DuplicateDeliveryWeekday(weekday uint8) *AppError {
	e := *ErrDuplicateDeliveryWeekday
	e.Code = ErrDuplicateDeliveryWeekday.Code
	e.FullCode = fmt.Sprintf("%s__%d", e.Code, weekday)

	// allow errors.Is
	e.err = ErrDuplicateDeliveryWeekday

	return &e
}
