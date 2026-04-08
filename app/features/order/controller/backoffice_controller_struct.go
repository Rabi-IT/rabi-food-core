package controller

import "github.com/Rabi-IT/rabi-food-core/features/order/usecases"

type OrderBackofficeController struct {
	usecase *usecases.OrderCase
}

func NewBackoffice(usecase *usecases.OrderCase) *OrderBackofficeController {
	return &OrderBackofficeController{usecase}
}
