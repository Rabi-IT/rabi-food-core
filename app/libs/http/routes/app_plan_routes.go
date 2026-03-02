package routes

import (
	"rabi-food-core/libs/http/controllers/app_plan_controller"

	"github.com/gofiber/fiber/v2"
)

func AppPlan(app *fiber.App, c *app_plan_controller.AppPlanController) {
	route := app.Group("/app_plan")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
}
