package payment_routes

import (
	payment_controller "github.com/Rabi-IT/rabi-food-core/features/payment/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

func PaymentProtected(app *fiber.App, c *payment_controller.PaymentController) {
	route := app.Group("/payment", middlewares.RequireTenantID)
	route.Get("/profile", c.GetProfile)
	route.Post("/intent", c.CreateIntent)
	route.Post("/confirm", c.Confirm)
}
