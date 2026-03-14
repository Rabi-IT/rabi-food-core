package docs

import (
	"github.com/gofiber/fiber/v2"
)

func ServeUI(c *fiber.Ctx) error {
	html := `
	<!doctype html>
	<html>
	<head>
		<title>Scalar API Reference</title>
		<meta charset="utf-8" />
		<meta
		name="viewport"
		content="width=device-width, initial-scale=1" />
		<link rel="icon" type="image/png" href="/docs/favicon.png" />
	</head>
	<body>
		<div id="app"></div>
		<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		<script>
		Scalar.createApiReference('#app', {
			url: '/docs/swagger.json',
		})
		</script>
	</body>
	</html>`

	c.Set("Content-Type", "text/html")

	return c.SendString(html)
}
