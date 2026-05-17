package payment_controller

import "github.com/Rabi-IT/rabi-food-core/features/payment/usecases"

type PaymentController struct {
	usecase *payment_usecases.PaymentCase
}

func New(usecase *payment_usecases.PaymentCase) *PaymentController {
	return &PaymentController{usecase}
}
