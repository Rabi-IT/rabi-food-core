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
		user := session.GetOriginalUser()

		b := logger.Get(uctx).
			With().
			Str("request_id", requestID).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Str("request_user", user)

		if session.TenantID != "" {
			b = b.Str("tenant_id", session.TenantID)
		}

		log := b.Logger()
		c.SetUserContext(logger.WithContext(uctx, log))
		return c.Next()
	}
}
