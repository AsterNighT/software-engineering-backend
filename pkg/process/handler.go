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
// @Tags patient
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {array} api.ReturnedData{data=[]Department}
// @Router /departments [GET]
func (h *RegistrationHandler) GetAllDepartments(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetDepartmentByID
// @Summary get a department by its ID
// @Tags patient
// @Description return a department's details by its ID
// @Param department_id path uint true "Department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Department}
// @Router /departments/{department_id} [GET]
func (h *RegistrationHandler) GetDepartmentByID(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateRegistrationByPatient
// @Summary patient register
// @Tags patient
// @Description return registration state
// @Param patient_id query uint true "Patient ID"
// @Param department_id query uint true "Department ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{data="register success"}
// @Router /patient/registration/create [POST]
func (h *RegistrationHandler) CreateRegistrationByPatient(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetRegistrationsByPatient
// @Summary get all registrations (patient view)
// @Tags patient
// @Description display all registrations of a patient
// @Param patient_id path uint true "Patient's ID"
// @Produce json
// @Success 200 {array} api.ReturnedData{data=[]Registration}
// @Router /patient/{patient_id}/registrations [GET]
func (h *RegistrationHandler) GetRegistrationsByPatient(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetRegistrationsByDoctor
// @Summary get all registrations (doctor view)
// @Tags doctor
// @Description display all registrations of a patient
// @Param doctor_id path uint true "Patient ID"
// @Produce json
// @Success 200 {array} api.ReturnedData{data=[]Registration}
// @Router /doctor/registrations/{doctor_id} [GET]
func (h *RegistrationHandler) GetRegistrationsByDoctor(c echo.Context) error {

	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// GetRegistrationDetailByPatient
// @Summary get a registration by its ID (patient view)
// @Tags patient
// @Description return a registration details by its ID
// @Param registration_id path uint true "Registration ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Registration}
// @Router /patient/registration/get/{registration_id} [GET]
func (h *RegistrationHandler) GetRegistrationDetailByPatient(c echo.Context) error {
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// this may not be used by doctor view's frontend
//// @Summary get a registration by its ID (doctor view)
//// @Tags doctor
//// @Description return a registration details by its ID
//// @Param registration_id
//// @Produce json
//// @Success 200 {object} api.ReturnedData{data=Registration}
//// @Router /doctor/registration/get/{registration_id} [GET]
//func (h *RegistrationHandler) GetRegistrationDetailByDoctor(c echo.Context) error {
//	c.Logger().Debug("hello world")
//	return c.JSON(http.StatusCreated, api.Return("ok", nil))
//}

// DeleteRegistrationByPatient
// @Summary delete a registration by its ID
// @Tags patient
// @Description return state
// @Param registration_id path uint true "Registration ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{"delete success"}
// @Router /patient/registration/delete/{registration_id} [GET]
func (h *RegistrationHandler) DeleteRegistrationByPatient(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// UpdateRegistrationStatus
// @Summary update registration status
// @Tags doctor
// @Description doctor change registration status
// @Param registration_id body uint true "Registration ID"
// @Param status body string true "Next registration's status"
// @Produce json
// @Success 200 {string} api.ReturnedData{"update success"}
// @Router /doctor/registration/update [POST]
func (h *RegistrationHandler) UpdateRegistrationStatus(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateRegistrationByDoctor
// @Summary create new registration (by doctor)
// @Tags doctor
// @Description the doctor create a new registration for the patient
// @Param patient_id body uint true "Patient ID"
// @Param department_id body uint true "Department ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{"create success"}
// @Router /doctor/registration/create [POST]
func (h *RegistrationHandler) CreateRegistrationByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateMileStoneByDoctor
// @Summary create milestone
// @Tags doctor
// @Description the doctor create milestone (type: array)
// @Param registration_id body uint true "Patient ID"
// @Param activity body string true "Milestone's activity"
// @Produce json
// @Success 200 {string} api.ReturnedData{"create success"}
// @Router /doctor/milestone/create [POST]
func (h *RegistrationHandler) CreateMileStoneByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// UpdateMileStoneByDoctor
// @Summary update milestone
// @Tags doctor
// @Description the doctor update milestone (check milestone)
// @Param milestone_id body uint "Milestone's ID"
// @Param checked body boolean true "Milestone is checked or not"
// @Produce json
// @Success 200 {string} api.ReturnedData{"update success"}
// @Router /doctor/milestone/update [POST]
func (h *RegistrationHandler) UpdateMileStoneByDoctor(c echo.Context) error {
	//verify identity
	c.Logger().Debug("hello world")
	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}


//// not included in our model
//// CreateOrder
//// @Summary create order
//// @Tags doctor
//// @Description the doctor create order for patient
//// @Param doctor_id, patient_id, medicines{}, description
//// @Produce json
//// @Success 200 {string} api.ReturnedData{"create order success"}
//// @Router /doctor/order/create [POST]
//func (h *RegistrationHandler) CreateOrder(c echo.Context) error {
//	//verify identity
//	c.Logger().Debug("hello world")
//	return c.JSON(http.StatusCreated, api.Return("ok", nil))
//}
