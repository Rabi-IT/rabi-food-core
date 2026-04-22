package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/subscription/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SubscriptionProtected registers subscription routes that require a valid JWT.
func SubscriptionProtected(app *fiber.App, c *controller.SubscriptionController) {
	route := app.Group("/subscription", middlewares.RequireTenantID)
	route.Post("/", c.Create)
	route.Get("/", c.Paginate)
	route.Get("/:id", c.GetByID)
	route.Put("/config", c.UpsertConfig)

	backoffice := app.Group("/backoffice/subscription", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
