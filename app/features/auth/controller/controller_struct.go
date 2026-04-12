package controller

import "github.com/Rabi-IT/rabi-food-core/features/auth/usecases"

type AuthController struct {
	usecase *usecases.AuthCase
}

func New(usecase *usecases.AuthCase) *AuthController {
	return &AuthController{usecase: usecase}
}
