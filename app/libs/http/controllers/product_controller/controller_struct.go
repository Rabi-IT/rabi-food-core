package product_controller

import (
	"rabi-food-core/usecases/product_case"
)

type ProductController struct {
	usecase *product_case.ProductCase
}

func New(usecase *product_case.ProductCase) *ProductController {
	return &ProductController{usecase}
}
