package app_plan_controller

import (
	"rabi-food-core/usecases/app_plan_case"
)

type AppPlanController struct {
	usecase *app_plan_case.AppPlanCase
}

func New(usecase *app_plan_case.AppPlanCase) *AppPlanController {
	return &AppPlanController{usecase}
}
