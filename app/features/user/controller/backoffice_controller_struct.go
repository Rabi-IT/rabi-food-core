package controller

import "github.com/Rabi-IT/rabi-food-core/features/user/usecases"

type UserBackofficeController struct {
	usecase *usecases.UserCase
}

func NewBackoffice(usecase *usecases.UserCase) *UserBackofficeController {
	return &UserBackofficeController{usecase}
}
