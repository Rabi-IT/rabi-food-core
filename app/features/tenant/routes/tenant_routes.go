package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/tenant/controller"

	"github.com/gofiber/fiber/v2"
)

func Tenant(app *fiber.App, c *controller.TenantController) {
	route := app.Group("/tenant")
	route.Get("/:id", c.GetByID)
	route.Patch("/:id", c.Patch)

	backoffice := app.Group("/backoffice/tenant")
	backoffice.Get("/", c.BackofficePaginate)
}
