package category_controller

import (
	"rabi-food-core/usecases/category_case"
)

type CategoryController struct {
	usecase *category_case.CategoryCase
}

func New(usecase *category_case.CategoryCase) *CategoryController {
	return &CategoryController{usecase}
}
