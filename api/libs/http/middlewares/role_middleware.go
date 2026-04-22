package middlewares

import (
	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gofiber/fiber/v2"
)

// RequireBackoffice is a middleware that checks if the user has backoffice role.
func RequireBackoffice(c *fiber.Ctx) error {
	if !app_context.GetSession(c.UserContext()).Role.IsBackoffice() {
		return errs.ErrForbidden
	}

	return c.Next()
}

// RequireTenant is a middleware that checks if the user has tenant role.
func RequireTenant(c *fiber.Ctx) error {
	if !app_context.GetSession(c.UserContext()).IsTenant() {
		return errs.ErrForbidden
	}

	return c.Next()
}
