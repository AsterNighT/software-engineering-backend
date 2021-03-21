package main

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()
	// How to use middleware
	app.Use(middleware.Logger())

	router.RegisterRouters(app)
	app.Logger.Fatal(app.Start(":12448"))
}
