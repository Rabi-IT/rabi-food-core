package user_controller

import (
	"net/http"
	g "rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// PatchUser godoc
// @Summary Patch user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce text/plain
// @Param id path string true "User ID"
// @Param user body g.PatchValues true "User patch data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [patch]
func (c *UserController) Patch(ctx *fiber.Ctx) error {
	data := g.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	updated, err := c.usecase.Patch(ctx.UserContext(), ctx.Params("id"), data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
