package usecases

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/golang-jwt/jwt/v5"
)

const scopedTokenTTL = 30 * time.Minute

type ExchangeTokenInput struct {
	TenantID string `json:"tenantId" validate:"required,uuid"`
}

type ExchangeTokenOutput struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

func (c *AuthCase) ExchangeToken(ctx context.Context, input *ExchangeTokenInput) (*ExchangeTokenOutput, error) {
	session := app_context.GetSession(ctx)

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

	return &ExchangeTokenOutput{
		AccessToken: signed,
		ExpiresIn:   int(scopedTokenTTL.Seconds()),
		TokenType:   "Bearer",
	}, nil
}
