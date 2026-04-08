package controller

import "github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"

type TenantBackofficeController struct {
	usecase *usecases.TenantCase
}

func NewBackoffice(usecase *usecases.TenantCase) *TenantBackofficeController {
	return &TenantBackofficeController{usecase}
}
