package middlewares

import (
	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gofiber/fiber/v2"
)

// RequireRole returns a middleware that rejects requests whose session role is not in the allowed list.
func RequireRole(roles ...auth.Role) fiber.Handler {
	allowed := make(map[auth.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		session := app_context.GetSession(c.UserContext())
		if _, ok := allowed[session.Role]; !ok {
			return errs.ErrForbidden
		}
		return c.Next()
	}
}
