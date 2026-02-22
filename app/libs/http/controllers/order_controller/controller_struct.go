package order_controller

import (
	"rabi-food-core/usecases/order_case"
)

type OrderController struct {
	usecase *order_case.OrderCase
}

func New(usecase *order_case.OrderCase) *OrderController {
	return &OrderController{usecase}
}
