package ui

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

type ScalarUI struct {
	api huma.API
}

func NewScalarUI(api huma.API) *ScalarUI {
	return &ScalarUI{api: api}
}

func (s *ScalarUI) RegisterRoutes(router fiber.Router) {
	docs := router.Group("/docs")
	docs.Get("/openapi.json", s.serveSchema)
	docs.Get("/", s.serveUI)
}

func (s *ScalarUI) serveSchema(c *fiber.Ctx) error {
	schema, err := json.MarshalIndent(s.api.OpenAPI(), "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to generate schema")
	}

	c.Set("Content-Type", "application/json")
	return c.Send(schema)
}

func (s *ScalarUI) serveUI(c *fiber.Ctx) error {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Products API — Docs</title>
    <style>body { margin: 0; }</style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/openapi.json"
      data-configuration='{"theme":"purple"}'
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
