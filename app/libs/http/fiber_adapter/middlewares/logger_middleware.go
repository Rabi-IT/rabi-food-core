package middlewares

import (
	"rabi-food-core/app_context"
	"rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uctx := c.UserContext()

		requestID, _ := c.Locals(requestid.ConfigDefault.ContextKey).(string)

		session := app_context.GetSession(uctx)
		userID, isBackoffice := session.GetOriginalUserID()

		l := logger.Get(uctx).
			With().
			Str(logger.RequestID, requestID).
			Str(logger.UserID, userID)

		if isBackoffice {
			l = l.Bool(logger.IsBackoffice, isBackoffice)
		}

		if session.TenantID != "" {
			l = l.Str(logger.TenantID, session.TenantID)
		}

		log := l.Logger()
		c.SetUserContext(logger.WithContext(uctx, log))

		log.Info().
			Str(logger.Path, c.Path()).
			Str(logger.Method, c.Method()).
			Msg("Incoming request")

		return c.Next()
	}
}
