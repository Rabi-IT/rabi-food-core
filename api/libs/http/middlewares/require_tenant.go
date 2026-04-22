package middlewares

import "github.com/gofiber/fiber/v2"

// RequireTenantID is a middleware that ensures the X-Tenant-ID header is present.
func RequireTenantID(c *fiber.Ctx) error {
	if c.Get("X-Tenant-ID") == "" {
		return fiber.NewError(fiber.StatusBadRequest, "MISSING_TENANT_ID")
	}

	return c.Next()
}
