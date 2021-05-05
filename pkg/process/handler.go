package process

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegistrationHandler struct {
}

// GetAllDepartments
// @Summary get all departments
// @Tags Process
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Department}
// @Router /departments [GET]
func (h *RegistrationHandler) GetAllDepartments(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetDepartmentByID
// @Summary get a department by its ID
// @Tags Process
// @Description return a department's details by its ID
// @Param DepartmentID path uint true "department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Department}
// @Router /department/{DepartmentID} [GET]
func (h *RegistrationHandler) GetDepartmentByID(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateRegistration
// @Summary create registration
// @Tags Process 
// @Description return registration state
// @Param PatientID body uint true "patient's ID"
// @Param DepartmentID body uint true "department ID"
// @Param RegistrationTime body string true "registration's time, basic format 'YYYY-MM-DD-<half-day>'"
// @Produce json
// @Success 200 {object} api.ReturnedData{}
// @Router /registration [POST]
func (h *RegistrationHandler) CreateRegistration(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// NOTE: we assume a doctor may help the patient commit registration

// GetRegistrationsByPatient
// @Summary get all registrations (patient view)
// @Tags Process 
// @Description display all registrations of a patient
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Registration}
// @Router /patient/registrations [GET]
func (h *RegistrationHandler) GetRegistrationsByPatient(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// NOTE: use cookie for identification

// GetRegistrationsByDoctor
// @Summary get all registrations (doctor view)
// @Tags Process 
// @Description display all registrations of a patient
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Registration}
// @Router /doctor/registrations [GET]
func (h *RegistrationHandler) GetRegistrationsByDoctor(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// NOTE: use cookie for identification

// GetRegistrationByPatient
// @Summary get a registration by its ID (patient view)
// @Tags Process 
// @Description return a registration details by its ID
// @Param RegistrationID path uint true "registration's ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Registration}
// @Router /patient/registration/{RegistrationID} [GET]
func (h *RegistrationHandler) GetRegistrationByPatient(c echo.Context) error {
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetRegistrationByDoctor
// @Summary get a registration by its ID (doctor view)
// @Tags Process
// @Description return a registration details by its ID
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Registration}
// @Router /doctor/registration/{RegistrationID} [GET]
func (h *RegistrationHandler) GetRegistrationByDoctor(c echo.Context) error {
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// UpdateRegistrationStatus
// @Summary update registration status
// @Tags Process
// @Description update registration status
// @Param RegistrationID path uint true "registration ID"
// @Param Status body string true "next status of current registration"
// @Param TerminatedCause body string false "if registration is shutdown, a cause is required"
// @Produce json
// @Success 200 {object} api.ReturnedData{}
// @Router /registration/{RegistrationID} [PUT]
func (h *RegistrationHandler) UpdateRegistrationStatus(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateMileStoneByDoctor
// @Summary create milestone
// @Tags Process
// @Description the doctor create milestone (type: array)
// @Param RegistrationID body uint true "registration's ID"
// @Param Activity body string true "milestone's activity"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone [POST]
func (h *RegistrationHandler) CreateMileStoneByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// UpdateMileStoneByDoctor
// @Summary update milestone
// @Tags Process
// @Description the doctor update milestone (check milestone)
// @Param MileStoneID path uint true "milestone's ID"
// @Param Activity body string false "updated milestone's activity"
// @Param Checked body boolean true "milestone is checked or not"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone/{MileStoneID} [PUT]
func (h *RegistrationHandler) UpdateMileStoneByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// DeleteMileStoneByDoctor
// @Summary delete milestone
// @Tags Process
// @Description the doctor delete milestone
// @Param MileStoneID path uint true "milestone's ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone/{MileStoneID} [DELETE]
func (h *RegistrationHandler) DeleteMileStoneByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// NOTE: use delete method, because there is no dependencies on MileStone
