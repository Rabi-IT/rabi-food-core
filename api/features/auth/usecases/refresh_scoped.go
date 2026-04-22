package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshScopedInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
	TenantID     string `json:"tenantId"     validate:"required,uuid"`
}

type RefreshScopedOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

func (c *AuthCase) RefreshScoped(ctx context.Context, input *RefreshScopedInput) (*RefreshScopedOutput, error) {
	out, err := c.gateway.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	session, err := parseGoTrueToken(out.AccessToken)
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	member, err := c.tenantCase.GetMembership(ctx, session.UserID, input.TenantID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errs.ErrForbidden
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   session.UserID,
		"email": session.Login,
		"iat":   now.Unix(),
		"exp":   now.Add(scopedTokenTTL).Unix(),
		"app_metadata": map[string]any{
			"role":        session.Role,
			"tenant_role": member.Role,
			"tenant_id":   input.TenantID,
		},
		"user_metadata": map[string]any{
			"name": session.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.AuthSecret))
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	return &RefreshScopedOutput{
		AccessToken:  signed,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    int(scopedTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// parseGoTrueToken extracts user session claims from a GoTrue-issued access token.
// Signature verification is skipped because the token comes directly from GoTrue (server-to-server).
// Tenant membership is still validated by the caller before generating any scoped token.
func parseGoTrueToken(tokenStr string) (app_context.UserSession, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return app_context.UserSession{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return app_context.UserSession{}, fmt.Errorf("invalid claims")
	}

	appMeta, _ := claims["app_metadata"].(map[string]any)
	userMeta, _ := claims["user_metadata"].(map[string]any)

	return app_context.UserSession{
		UserID: fmt.Sprint(claims["sub"]),
		Login:  fmt.Sprint(claims["email"]),
		Role:   auth.Role(claimString(appMeta, "role")),
		Name:   claimString(userMeta, "name"),
	}, nil
}

func claimString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	return fmt.Sprint(v)
}
