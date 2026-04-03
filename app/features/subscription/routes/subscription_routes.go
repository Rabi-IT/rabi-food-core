package routes

import (
	"rabi-food-core/features/subscription/controller"

	"github.com/gofiber/fiber/v2"
)

func Subscription(app *fiber.App, c *controller.SubscriptionController) {
	route := app.Group("/subscription")
	route.Post("/", c.Create)
	route.Get("/:id", c.GetByID)
	route.Put("/config", c.UpsertConfig)
}
