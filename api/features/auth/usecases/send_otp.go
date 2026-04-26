package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"crypto/rand"
	"encoding/hex"
)

type SendOTPInput struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name"  validate:"required"`
}

func (c *AuthCase) SendOTP(ctx context.Context, input *SendOTPInput) error {
	_, err := c.gateway.CreateUser(ctx, g.CreateUserInput{
		Email:    input.Email,
		Name:     input.Name,
		Password: randomPassword(),
		Role:     auth.User,
	})
	if err != nil && !errors.Is(err, errs.ErrUserAlreadyExists) {
		return fmt.Errorf("send-otp create user: %w", err)
	}

	return c.gateway.SendOTP(ctx, g.SendOTPInput{
		Email: input.Email,
		Name:  input.Name,
	})
}

func randomPassword() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
