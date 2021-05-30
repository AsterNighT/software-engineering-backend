package main

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/database"
	"github.com/AsterNighT/software-engineering-backend/pkg/router"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// initialize database
	db := database.InitDB()

	// create echo instance
	app := echo.New()
	app.Validator = &utils.CustomValidator{Validator: validator.New()}
	app.Use(middleware.Logger())

	app.Use(database.ContextDB(db))

	err := router.RegisterRouters(app)
	if err != nil {
		panic(err)
	}
	// database.InitDb()
	app.Logger.Fatal(app.Start(":12448"))
}
