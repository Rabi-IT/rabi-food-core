package controller

import (
	"net/http"

	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	parser "github.com/Rabi-IT/rabi-food-core/libs/http/parser"

	"github.com/gofiber/fiber/v2"
)

// GetByID godoc
// @Summary Get user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body g.GetByIDFilter false "Optional user filter data"
// @Success 200 {object} g.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [get].
func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	filter := g.GetByIDFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return err
	}

	filter.ID = ctx.Params("id")
	data, err := c.usecase.GetByID(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
