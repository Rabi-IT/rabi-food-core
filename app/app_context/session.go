package app_context

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
)

type sessionKey string

const SessionKey sessionKey = "session"

type UserSession struct {
	UserID string
	Name   string
	Login  string
	Role   auth.Role
}

func GetSession(ctx context.Context) UserSession {
	session, ok := ctx.Value(SessionKey).(*UserSession)
	if !ok {
		return UserSession{}
	}

	return *session
}

func WithSession(ctx context.Context, s *UserSession) context.Context {
	return context.WithValue(ctx, SessionKey, s)
}

