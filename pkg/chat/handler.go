package chat

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
}

type KeywordHandler struct {
}

// @Summary New a chat
// @Description
// @Tags Chat
// @Produce json
// @Param chatDetail body Chat true "patient ID, doctor ID and other chat details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/chat [POST]
func (h *ChatHandler) NewChat(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Delete a chat
// @Description
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/chat/{chatID} [DELETE]
func (h *ChatHandler) DeleteChatByChatID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get the lastest message
// @Description
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Success 200 {object} api.ReturnedData{data=Message}
// @Router /patient/{patientID}/chat/{chatID} [GET]
func (h *ChatHandler) GetLastMessage(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary New a Message
// @Description
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Param messageDetail body Message true "chat ID, type and other message details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/chat/{chatID}/message [POST]
func (h *ChatHandler) NewMessage(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Delete a message by message id
// @Description Can be viewed as recall a message
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Param messageID path uint true "message ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/chat/{chatID}/message/{messageID} [DELETE]
func (h *ChatHandler) DeleteMessageByMessageID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get a message by message ID
// @Description
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Param messageID path uint true "message ID"
// @Success 200 {object} api.ReturnedData{data=Message}
// @Router /patient/{patientID}/chat/{chatID}/message/{messageID} [GET]
func (h *ChatHandler) GetMessageByMessageID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get a message list by chat ID
// @Description
// @Tags Chat
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param chatID path uint true "chat ID"
// @Success 200 {object} api.ReturnedData{data=[]Message}
// @Router /patient/{patientID}/chat/{chatID} [GET]
func (h *ChatHandler) GetMessagesByChatID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get a keyword list by message id
// @Description Parse the message and get valid keywords
// @Tags Chat
// @Produce json
// @Param messageID path uint true "message ID"
// @Success 200 {object} api.ReturnedData{[]Keyword}
// @Router /patient/{patientID}/chat/{chatID}/message/{messageID}  [GET]
func (h *ChatHandler) GetKeywordsByMessageID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get a catagory list by keyword id
// @Description
// @Tags Chat
// @Produce json
// @Param keywordID path uint true "keyword ID"
// @Success 200 {object} api.ReturnedData{[]Catagory}
// @Router /keyword/{keywordID} [GET]
func (h *KeywordHandler) GetCatagorysByKeywordID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get questions by catagory id
// @Description
// @Tags Chat
// @Produce json
// @Param catagoryID path uint true "catagory ID"
// @Success 200 {object} api.ReturnedData{data=[]string}
// @Router /keyword/{keywordID}/catagory/{catagoryID}  [GET]
func (h *KeywordHandler) GetQuestionsByCatagoryID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("ok", nil))
}
