package process

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegistrationHandler struct {
}

// @Summary get all departments
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {array} Department
// @Router /departments [GET]
func (h *RegistrationHandler) GetAllDepartments(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("pong", nil))
}
