package main

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/pkg/router"
	"github.com/labstack/echo/v4"
)

var Log echo.Logger

func main() {
	app := echo.New()
	Log = app.Logger
	Log.SetLevel(1)
	router.RegisterRoutersRegisterRouters(app)
	Log.Fatal(app.Start(":12448"))
}
