package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
)

type OrderController struct {
	usecase *usecases.OrderCase
}

func New(usecase *usecases.OrderCase) *OrderController {
	return &OrderController{usecase}
}
