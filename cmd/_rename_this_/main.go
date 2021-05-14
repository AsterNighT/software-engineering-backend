package main

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()
	app.Use(middleware.Logger())

	router.RegisterRouters(app)
	// database.InitDb()
	app.Logger.Fatal(app.Start(":12448"))
}
