package main

import (
	"log"
	"net/http"

	"github.com/AsterNighT/software-engineering-backend/pkg/database"
	"github.com/AsterNighT/software-engineering-backend/pkg/router"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	// load envs
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// initialize database
	db := database.InitDB()

	// create echo instance
	app := echo.New()
	app.Validator = &utils.CustomValidator{Validator: validator.New()}
	app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Cookies", "authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))
	app.Use(database.ContextDB(db))

	err := router.RegisterRouters(app)
	if err != nil {
		panic(err)
	}
	// database.InitDb()
	app.Logger.Fatal(app.Start(":12448"))
}
