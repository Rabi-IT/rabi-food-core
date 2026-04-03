package controller

import (
	"rabi-food-core/features/category/usecases"
)

type CategoryController struct {
	usecase *usecases.CategoryCase
}

func New(usecase *usecases.CategoryCase) *CategoryController {
	return &CategoryController{usecase}
}
