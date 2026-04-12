package routes

import (
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/features/tenant/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Tenant(app *fiber.App, c *controller.TenantController) {
	route := app.Group("/tenant")
	route.Get("/:id", c.GetByID)
	route.Patch("/:id", c.Patch)
	route.Post("/:id/customers", c.RegisterCustomer)

	backoffice := app.Group("/backoffice/tenant", middlewares.RequireRole(auth.Backoffice))
	backoffice.Get("/", c.BackofficePaginate)
}
