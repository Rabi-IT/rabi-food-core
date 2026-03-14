package user_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/user_case"

	"github.com/gofiber/fiber/v2"
)

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce text/plain
// @Param user body user_case.CreateInput true "User data"
// @Success 201 {string} string "Created user ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /user/ [post]
func (c *UserController) Create(ctx *fiber.Ctx) error {
	data := &user_case.CreateInput{}
	err := ctx.BodyParser(data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	id, err := c.usecase.Create(ctx.UserContext(), data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
