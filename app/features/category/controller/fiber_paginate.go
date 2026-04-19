package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate categories
// @Description Paginate categories using query filters
// @Tags categories
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param name query string false "Category name"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /category/ [get].
func (c *CategoryController) Paginate(ctx *fiber.Ctx) error {
	filter := gateway.PaginateFilter{}
	if err := ctx.QueryParser(&filter); err != nil {
		return err
	}
	filter.TenantID = ctx.Get("X-Tenant-ID")

	paginate := database.PaginateInput{
		Page:     ctx.QueryInt("Page", database.DefaultPage),
		PageSize: ctx.QueryInt("PageSize", database.DefaultPageSize),
	}

	result, err := c.usecase.Paginate(ctx.UserContext(), filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
