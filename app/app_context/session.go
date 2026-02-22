package app_context

import (
	"context"
	"rabi-food-core/domain"
)

type sessionKey string

const SessionKey sessionKey = "session"

type UserSession struct {
	UserID         string
	TenantID       string
	Name           string
	Login          string
	OriginalUserID string
	Role           domain.Role
}

func (u *UserSession) GetOriginalUserID() (string, bool) {
	if u.Role.IsBackoffice() {
		return u.OriginalUserID, true
	}

	return u.UserID, false
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
