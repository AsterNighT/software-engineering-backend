package cases

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
)

type CaseHandler struct {
}

// @Summary Get case by case id
// @Description
// @Produce json
// @Param  id path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/{id} [GET]
func (h *CaseHandler) GetCaseByCaseID(c echo.Context) error {
	// what is in here does not matter for now, just filling the swaggo comments and name the function
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("pong", nil))
}
