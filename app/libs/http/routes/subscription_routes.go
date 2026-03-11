package routes

import (
	"rabi-food-core/libs/http/controllers/subscription_controller"

	"github.com/gofiber/fiber/v2"
)

func Subscription(app *fiber.App, c *subscription_controller.SubscriptionController) {
	route := app.Group("/subscription")
	route.Post("/", c.Create)
	route.Get("/:id", c.GetByID)
}
