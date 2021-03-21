package router

import (
	"net/http"

	_ "github.com/AsterNighT/software-engineering-backend/docs" // swagger doc
	"github.com/labstack/echo/v4"
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

// @host localhost:12448
// @BasePath /api
func RegisterRouters(app *echo.Echo) error {

	app.GET("/swagger/*", echoSwagger.WrapHandler)
	{
		router := app.Group("/api")
		router.GET("/ping", pingHandler)
		{
			// Use nested scopes and shadowing for subgroups
			// router = router.Group("/case")
		}

	}
	return nil
}

// @Summary Test server up statue
// @Description respond to a ping request from client
// @Produce plain
// @Success 200 {string} string	"Good, server is up"
// @Router /ping [GET]
func pingHandler(c echo.Context) error {
	c.Logger().Debug("hello world")
	return c.String(http.StatusOK, "pong")
}
