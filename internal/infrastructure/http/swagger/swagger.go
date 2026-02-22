package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Honda Leasing API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger/doc.yaml",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`

// RegisterRoutes sets up Swagger UI routes on the given Gin engine.
// Access: GET /swagger/index.html
func RegisterRoutes(router *gin.Engine) {
	swagger := router.Group("/swagger")
	{
		swagger.GET("/index.html", func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
		})

		swagger.GET("/doc.yaml", func(c *gin.Context) {
			c.File("api/openapi/swagger.yaml")
		})
	}
}
