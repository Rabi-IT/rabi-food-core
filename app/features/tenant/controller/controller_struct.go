package controller

import "github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"

type TenantController struct {
	usecase *usecases.TenantCase
}

func New(usecase *usecases.TenantCase) *TenantController {
	return &TenantController{usecase}
}
