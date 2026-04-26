package middlewares

import (
	"errors"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	lib "github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationErrorResponse struct {
	Code   string                `json:"code"`
	Errors []lib.ValidationError `json:"errors"`
}

var (
	errBodyInternalError = []byte(`{"code":"INTERNAL_ERROR"}`)
	errBodyInvalidToken  = []byte(`{"code":"INVALID_TOKEN"}`)
)

// ErrorHandler is a middleware that handles errors occurring during request processing.
// Logging is handled by the Canonical Log Line middleware — this function only formats the response.
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	//nolint:errorlint
	validationErr, ok := err.(validator.ValidationErrors)
	if ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
			Code:   "VALIDATION_ERROR",
			Errors: lib.ParseValidationError(validationErr),
		})
	}

	if apiErr, ok := errors.AsType[*errs.ApiError](err); ok {
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	if e, ok := errors.AsType[*fiber.Error](err); ok {
		if e.Code == fiber.StatusUnauthorized {
			return ctx.Status(fiber.StatusUnauthorized).Type("json").Send(errBodyInvalidToken)
		}
		return ctx.Status(e.Code).Type("json").Send(errBodyInternalError)
	}

	return ctx.Status(http.StatusInternalServerError).Type("json").Send(errBodyInternalError)
}
