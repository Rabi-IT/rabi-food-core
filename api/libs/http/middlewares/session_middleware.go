package middlewares

import (
	"errors"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidToken  = errors.New("INVALID_TOKEN")
	errInvalidClaims = errors.New("INVALID_CLAIMS")
)

// Session is a middleware that extracts user session information from a GoTrue JWT token.
func Session(c *fiber.Ctx) error {
	token, ok := c.Context().UserValue("user").(*jwt.Token)
	if !ok || !token.Valid {
		return errInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errInvalidClaims
	}

	appMeta, _ := claims["app_metadata"].(map[string]any)
	userMeta, _ := claims["user_metadata"].(map[string]any)

	session := &app_context.UserSession{
		UserID:     fmt.Sprint(claims["sub"]),
		TenantID:   stringFromMap(appMeta, "tenant_id"),
		TenantRole: tenant.Role(stringFromMap(appMeta, "tenant_role")),
		Login:      fmt.Sprint(claims["email"]),
		Name:       stringFromMap(userMeta, "name"),
		Role:       auth.Role(stringFromMap(appMeta, "role")),
	}

	uctx := c.UserContext()
	logger.GetWideEvent(uctx).ActorID = session.UserID

	c.SetUserContext(app_context.WithSession(uctx, session))

	return c.Next()
}

func stringFromMap(m map[string]any, key string) string {
	if m == nil {
		return ""
	}

	v, ok := m[key]
	if !ok {
		return ""
	}

	return fmt.Sprint(v)
}
