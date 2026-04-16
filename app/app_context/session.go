package app_context

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
)

type sessionKey string

const sessionCtxKey sessionKey = "session"

type UserSession struct {
	UserID string
	// The tenant this user manages.
	// Empty if the user has no tenant (e.g. User, Backoffice).
	TenantID string
	Name     string
	Login    string
	Role     auth.Role
}

func GetSession(ctx context.Context) UserSession {
	session, ok := ctx.Value(sessionCtxKey).(*UserSession)
	if !ok {
		return UserSession{}
	}
	return *session
}

func WithSession(ctx context.Context, s *UserSession) context.Context {
	return context.WithValue(ctx, sessionCtxKey, s)
}
