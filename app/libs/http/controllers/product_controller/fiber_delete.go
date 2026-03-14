package product_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Delete godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Produce text/plain
// @Param id path string true "Product ID"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /product/{id} [delete].
func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	deleted, err := c.usecase.Delete(ctx.UserContext(), id)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
