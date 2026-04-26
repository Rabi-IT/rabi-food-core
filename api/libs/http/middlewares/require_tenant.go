package middlewares

import (
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/gofiber/fiber/v2"
)

// RequireTenantID is a middleware that ensures the X-Tenant-ID header is present.
func RequireTenantID(c *fiber.Ctx) error {
	if c.Get("X-Tenant-ID") == "" {
		return errs.ErrMissingTenantID
	}

	return c.Next()
}
