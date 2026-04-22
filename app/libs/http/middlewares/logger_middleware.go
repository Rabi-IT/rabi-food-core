package middlewares

import (
	"errors"
	"time"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
)

func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		uctx := c.UserContext()

		// Build the request-scoped logger with identity fields.
		requestID, _ := c.Locals(requestid.ConfigDefault.ContextKey).(string)

		wd := logger.Acquire()
		defer logger.Release(wd)

		wd.Method = c.Method()
		wd.Path = c.Path()
		wd.RequestID = requestID
		query := c.Context().QueryArgs().String()
		if query != "" {
			wd.Query = query
		}

		c.SetUserContext(logger.WithWideEvent(uctx, wd))

		// Execute the handler chain.
		err := c.Next()
		// After this point, changing the WideEvent is not thread-safe, so we use mutex.

		// Determine status code and log level.
		statusCode := c.Response().StatusCode()
		if err != nil {
			statusCode = classifyStatus(err)
		}

		isValidationError := err != nil && errors.As(err, new(validator.ValidationErrors))
		errToLog := err
		if isValidationError {
			// Don't log ValidationErrors as errors to avoid noise and high cardinality in log aggregation
			errToLog = nil
		}
		wd.FinishRequest(statusCode, errToLog, time.Since(start).Milliseconds())

		l := logger.New().Logger()
		var event *zerolog.Event
		switch {
		case statusCode < 400:
			event = l.Info()
		case statusCode >= 500:
			event = l.Error()
		case isValidationError:
			event = l.Info()
		case statusCode >= 400:
			event = l.Warn()
		default:
			event = l.Info()
		}

		wd.MarshalZerologObject(event)
		event.Send()

		return err
	}
}

// classifyStatus mirrors the ErrorHandler logic to determine the HTTP status
// that will be sent to the client, without allocating or mutating the response.
func classifyStatus(err error) int {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return fiber.StatusBadRequest
	}

	var apiErr *errs.ApiError
	if errors.As(err, &apiErr) {
		return apiErr.Status
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return fiber.StatusInternalServerError
}
