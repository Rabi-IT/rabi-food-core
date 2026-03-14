package category_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} category_gateway.GetByIDOutput
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /category/{id} [get]
func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.usecase.GetByID(ctx.UserContext(), id)

	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
