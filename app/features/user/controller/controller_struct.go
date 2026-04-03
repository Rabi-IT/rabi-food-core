package controller

import "github.com/Rabi-IT/rabi-food-core/features/user/usecases"

type UserController struct {
	usecase *usecases.UserCase
}

func New(usecase *usecases.UserCase) *UserController {
	return &UserController{usecase}
}
