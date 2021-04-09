package router

import (
	"net/http"

	"github.com/AsterNighT/software-engineering-backend/api"
	_ "github.com/AsterNighT/software-engineering-backend/docs" // swagger doc
	"github.com/AsterNighT/software-engineering-backend/pkg/chat"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type BasicHandler struct {
}

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

	var h BasicHandler
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.GET("/swagger", h.RedirectToSwagger)
	{
		router := app.Group("/api")
		router.GET("/ping", h.Ping)
		{
			// Use nested scopes and shadowing for subgroups
			var h chat.ChatHandler
			router = router.Group("/patient")
			router.POST("/:patientID/chat", h.NewChat)
			router.DELETE("/:patientID/chat/:chatID", h.DeleteChatByChatID)
			router.GET("/:patientID/chat/:chatID", h.GetLastMessage)
			router.POST("/:patientID/chat/:chatID/message", h.NewMessage)
			router.DELETE("/:patientID/chat/:chatID/message/:messageID", h.DeleteMessageByMessageID)
			router.GET("/:patientID/chat/:chatID/message/:messageID", h.GetMessageByMessageID)
			router.GET("/:patientID/chat/:chatID", h.GetMessagesByChatID)
			router.GET("/:patientID/chat/:chatID/message/:messageID", h.GetKeywordsByMessageID)
		}
		{
			var h chat.KeywordHandler
			router.GET("/keyword/:keywordID", h.GetCatagorysByKeywordID)
		}
		{
			var h chat.CatagoryHandler
			router.GET("/catagory/:catagoryID", h.GetQuestionsByCatagoryID)
		}
	}
	return nil
}

// @Summary Test server up statue
// @Description respond to a ping request from client
// @Produce json
// @Success 200 {object} api.ReturnedData	"Good, server is up"
// @Router /ping [GET]
func (h *BasicHandler) Ping(c echo.Context) error {
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("pong", nil))
}

func (h *BasicHandler) RedirectToSwagger(c echo.Context) error {
	c.Response().Header().Set("Location", "swagger/index.html")
	return c.NoContent(http.StatusMovedPermanently)
}
