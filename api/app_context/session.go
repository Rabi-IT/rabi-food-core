package app_context

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
)

type sessionKey string

const sessionCtxKey sessionKey = "session"

type UserSession struct {
	UserID     string
	TenantID   string
	TenantRole tenant.Role
	Name       string
	Login      string
	Role       auth.Role
}

func (s UserSession) IsTenant() bool {
	return s.TenantRole != "" && s.TenantID != ""
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
