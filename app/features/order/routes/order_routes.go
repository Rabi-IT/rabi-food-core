package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/order/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Order(app *fiber.App, c *controller.OrderController) {
	route := app.Group("/order")
	route.Post("/", middlewares.RequireTenantID, c.Create)
	route.Delete("/:id", middlewares.RequireTenantID, c.Delete)
	route.Get("/:id", middlewares.RequireTenantID, c.GetByID)
	route.Get("/", middlewares.RequireTenantID, c.Paginate)
	route.Post("/:id/payments/confirm", c.ConfirmPayment)

	backoffice := app.Group("/backoffice/order", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
