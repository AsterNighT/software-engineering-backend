package main

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/database"
	"github.com/AsterNighT/software-engineering-backend/pkg/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// initialize database
	database.InitDb()

	// create echo instance
	app := echo.New()
	app.Use(middleware.Logger())

	err := router.RegisterRouters(app)
	if err != nil {
		panic(err)
	}
	// database.InitDb()
	app.Logger.Fatal(app.Start(":12448"))
}
