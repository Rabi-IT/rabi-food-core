package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	parser "github.com/Rabi-IT/rabi-food-core/libs/http/parser"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

// Delete godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce text/plain
// @Param id path string true "Category ID"
// @Param category body gateway.DeleteFilter false "Optional delete filter data"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /category/{id} [delete].
func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	filter := gateway.DeleteFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return ctx.JSON(err)
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "delete-category"

	filter.ID = ctx.Params("id")
	filter.TenantID = app_context.GetSession(uctx).TenantID
	deleted, err := c.usecase.Delete(uctx, filter)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
