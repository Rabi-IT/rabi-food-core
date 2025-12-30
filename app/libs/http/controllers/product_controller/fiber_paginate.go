package product_controller

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (c *ProductController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := product_gateway.PaginateFilter{}
	err = ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(ctx.UserContext(), filter, paginate)

	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
