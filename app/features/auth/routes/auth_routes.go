package routes

import (
	domain_auth "github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/features/auth/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Auth registers public auth routes (before JWT middleware).
func Auth(app *fiber.App, c *controller.AuthController) {
	route := app.Group("/auth")
	route.Post("/signup", c.SignUp)
	route.Post("/signin", c.SignIn)
	route.Post("/refresh", c.Refresh)
}

// AuthProtected registers auth routes that require a valid JWT.
func AuthProtected(app *fiber.App, c *controller.AuthController) {
	authGroup := app.Group("/auth")
	authGroup.Post("/signout", c.SignOut)

	user := app.Group("/user")
	user.Get("/me", c.GetMe)
	user.Get("/:id", c.GetByID)
	user.Patch("/:id", c.Patch)
	user.Delete("/:id", c.Delete)
	user.Get("/", c.Paginate)

	backoffice := app.Group("/backoffice/user", middlewares.RequireRole(domain_auth.Backoffice))
	backoffice.Get("/", c.Paginate)
}
