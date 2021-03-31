package process

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegistrationHandler struct {
}

// @Summary get all departments
// @Tags departments
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {array} api.ReturnedData{data=[]Department}
// @Router /departments [GET]
func (h *RegistrationHandler) GetAllDepartments(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("pong", nil))
}



// @Summary get all departments
// @Tags departments
// @Description display all departments of a hospital
// @Param department_id path uint true "Department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Department}
// @Router /departments [GET]
func (h *RegistrationHandler) GetDepartment(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("pong", nil))
}
