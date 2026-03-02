package app_plan_case

import g "rabi-food-core/libs/database/gateways/app_plan_gateway"

// AppPlanCase encapsulates the business logic related to app plans.
type AppPlanCase struct {
	gateway g.AppPlanGateway
}

// New creates a new instance of AppPlanCase with the provided AppPlanGateway.
func New(gateway g.AppPlanGateway) *AppPlanCase {
	return &AppPlanCase{gateway}
}
