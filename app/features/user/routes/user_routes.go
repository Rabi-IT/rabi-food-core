package routes

import (
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/features/user/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App, c *controller.UserController) {
	route := app.Group("/user")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)

	backoffice := app.Group("/backoffice/user", middlewares.RequireRole(auth.Backoffice))
	backoffice.Get("/", c.BackofficePaginate)
}
