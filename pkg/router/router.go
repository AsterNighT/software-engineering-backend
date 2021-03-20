package router

import (
	"net/http"

	"github.com/labstack/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /api
func RegisterRouters(app *echo.Echo) error {
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.GET("/ping", pingHandler)
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
