package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/gofiber/fiber/v2"
)

// Delete godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce text/plain
// @Param id path string true "User ID"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [delete].
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "delete-user"
	deleted, err := c.usecase.Delete(uctx, id)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
