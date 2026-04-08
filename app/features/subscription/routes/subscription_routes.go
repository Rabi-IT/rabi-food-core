package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/subscription/controller"

	"github.com/gofiber/fiber/v2"
)

func Subscription(app *fiber.App, c *controller.SubscriptionController) {
	route := app.Group("/subscription")
	route.Post("/", c.Create)
	route.Get("/", c.Paginate)
	route.Get("/:id", c.GetByID)
	route.Put("/config", c.UpsertConfig)

	backoffice := app.Group("/backoffice/subscription")
	backoffice.Get("/", c.BackofficePaginate)
}
