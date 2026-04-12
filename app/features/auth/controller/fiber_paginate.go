package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "10"))
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "user-paginate"

	out, err := c.usecase.Paginate(uctx, page, pageSize)
	if err != nil {
		return err
	}

	return ctx.JSON(out)
}
