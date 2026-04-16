package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/tenant/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Tenant(app *fiber.App, c *controller.TenantController) {
	app.Post("/tenant/:id/enroll-customers", c.EnrollCustomer)

	route := app.Group("/tenant", middlewares.RequireTenant)
	route.Get("/me", c.GetMe)
	route.Patch("/me", c.PatchMe)

	backoffice := app.Group("/backoffice/tenant", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
	backoffice.Get("/:id", c.GetByID)
}
