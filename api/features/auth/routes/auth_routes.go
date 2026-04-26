package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/auth/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Auth registers public auth routes (before JWT middleware).
func Auth(app *fiber.App, c *controller.AuthController) {
	route := app.Group("/auth")
	route.Post("/signup", c.SignUp)
	route.Post("/signin", c.SignIn)
	route.Post("/otp", c.SendOTP)
	route.Post("/otp/verify", c.VerifyOTP)
	route.Post("/refresh", c.Refresh)
	route.Post("/refresh-scoped", c.RefreshScoped)
}

// AuthProtected registers auth routes that require a valid JWT.
func AuthProtected(app *fiber.App, c *controller.AuthController) {
	authGroup := app.Group("/auth")
	authGroup.Post("/signout", c.SignOut)
	authGroup.Post("/token/exchange", c.ExchangeToken)

	user := app.Group("/user")
	user.Get("/me", c.GetMe)
	user.Patch("/me", c.Patch)

	backoffice := app.Group("/backoffice/user", middlewares.RequireBackoffice)
	backoffice.Get("/", c.Paginate)
}
