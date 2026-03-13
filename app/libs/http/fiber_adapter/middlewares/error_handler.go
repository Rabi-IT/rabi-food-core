package middlewares

import (
	"errors"
	"net/http"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"
	lib "rabi-food-core/libs/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidationErrorResponse represents the structure of validation error responses.
// It is used to return detailed information about validation errors to the client.
type ValidationErrorResponse struct {
	Errors []lib.ValidationError `json:"errors"`
}

// ErrorHandler is a middleware that handles errors occurring during request processing.
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	//nolint:errorlint
	validationErr, ok := err.(validator.ValidationErrors)
	if ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
			Errors: lib.ParseValidationError(validationErr),
		})
	}

	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		return ctx.Status(appErr.Status).JSON(appErr)
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return ctx.Status(e.Code).SendString(err.Error())
	}

	logger.Get(ctx.UserContext()).Error().Err(err).Msg("internal server error")

	return ctx.Status(http.StatusInternalServerError).SendString("internal server error")
}
