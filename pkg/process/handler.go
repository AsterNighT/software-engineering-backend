package process

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)



type ProcessHandler struct {
}

// GetAllDepartments
// @Summary get all departments
// @Tags Process
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Department}
// @Router /departments [GET]
func (h *ProcessHandler) GetAllDepartments(c echo.Context) error {
	DB := utils.GetDB()

	var Departments []Department
	DB.Find(&Departments)

	return c.JSON(http.StatusOK, api.Return("ok", Departments))
}

// GetDepartmentByID
// @Summary get a department by its ID
// @Tags Process
// @Description return a department's details by its ID
// @Param DepartmentID path uint true "department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Department}
// @Router /department/{DepartmentID} [GET]
func (h *ProcessHandler) GetDepartmentByID(c echo.Context) error {

	db := utils.GetDB()
	var department Department
	db.First(&department, c.Param("departmentID"))
	c.Logger().Debug("GetDepartmentByID")
	return c.JSON(http.StatusOK, api.Return("ok", department))

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
func (h *ProcessHandler) CreateRegistration(c echo.Context) error {
	registration := Registration{
		//TODO: need to allocate an doctor for patient by algorithm
		DoctorID:     uint(c.Get("DoctorID").(int)),
		PatientID:    uint(c.Get("PatientID").(int)),
		DepartmentID: uint(c.Get("DepartmentID").(int)),
		Date:         time.Now(),
	}
	db := utils.GetDB()
	if err := db.Create(registration).Error; err != nil {
		c.Logger().Error("Registration insert failed!")
	}
	c.Logger().Debug("CreateRegistration")
	return c.JSON(http.StatusOK, api.Return("ok", nil))
}

// NOTE: we assume a doctor may help the patient commit registration

// GetRegistrationsByPatient
// @Summary get all registrations (patient view)
// @Tags Process
// @Description display all registrations of a patient
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Registration}
// @Router /patient/registrations [GET]
func (h *ProcessHandler) GetRegistrationsByPatient(c echo.Context) error {
	patientID, err := strconv.ParseUint(c.Param("patientID"), 10, 64)
	if err != nil {
		c.Logger().Error("Patient not found!")
	}
	db := utils.GetDB()
	var registrations Registration
	db.Where("PatientID = ?", patientID).Find(&registrations)
	c.Logger().Debug("GetRegistrationsByPatient")
	return c.JSON(http.StatusOK, api.Return("ok", registrations))
}

// NOTE: use cookie for identification

// GetRegistrationsByDoctor
// @Summary get all registrations (doctor view)
// @Tags Process
// @Description display all registrations of a patient
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Registration}
// @Router /doctor/registrations [GET]
func (h *ProcessHandler) GetRegistrationsByDoctor(c echo.Context) error {
	doctorID, err := strconv.ParseUint(c.Param("doctorID"), 10, 64)
	if err != nil {
		c.Logger().Error("Doctor not found!")
	}
	db := utils.GetDB()
	var registrations Registration
	db.Where("DoctorID = ?", doctorID).Find(&registrations)
	c.Logger().Debug("GetRegistrationsByDoctor")
	return c.JSON(http.StatusCreated, api.Return("ok", registrations))
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
func (h *ProcessHandler) GetRegistrationByPatient(c echo.Context) error {
	//Must make sure RegistrationID in patientID's registration list.
	patientID, err := strconv.ParseUint(c.Param("patientID"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	registrationID, err := strconv.ParseUint(c.Param("registrationID"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	db := utils.GetDB()
	var registration Registration
	db.Where("PatientID = ? and RegistrationID = ?", patientID, registrationID).First(&registration)
	c.Logger().Debug("GetRegistrationByPatient")
	return c.JSON(http.StatusCreated, api.Return("ok", registration))
}

// GetRegistrationByDoctor
// @Summary get a registration by its ID (doctor view)
// @Tags Process
// @Description return a registration details by its ID
// @Produce json
// @Success 200 {object} api.ReturnedData{data=Registration}
// @Router /doctor/registration/{RegistrationID} [GET]
func (h *ProcessHandler) GetRegistrationByDoctor(c echo.Context) error {
	//Must make sure RegistrationID in patientID's registration list.
	doctorID, err := strconv.ParseUint(c.Param("doctorID"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	registrationID, err := strconv.ParseUint(c.Param("registrationID"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	db := utils.GetDB()
	var registration Registration
	db.Where("DoctorID = ? and RegistrationID = ?", doctorID, registrationID).First(&registration)
	c.Logger().Debug("GetRegistrationByDoctor")
	return c.JSON(http.StatusCreated, api.Return("ok", registration))
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
func (h *ProcessHandler) UpdateRegistrationStatus(c echo.Context) error {
	// TODO verify identity
	DB := utils.GetDB()

	RegistrationID, err := strconv.ParseUint(c.QueryParam("RegistrationID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	// TODO validation for status
	Status := RegistrationStatusEnum(c.QueryParam("Status"))
	TerminatedCause := c.QueryParam("TerminatedCause")

	var Registration Registration
	err = DB.First(&Registration, RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	Registration.Status = Status

	if Status == terminated {
		if TerminatedCause != "" {
			Registration.Status = Status
			Registration.TerminatedCause = TerminatedCause
		} else {
			return c.JSON(http.StatusBadRequest, api.Return("ok", "Missing terminated causes"))
		}
	} else {
		Registration.Status = Status
	}

	return c.JSON(http.StatusCreated, api.Return("ok", nil))
}

// CreateMileStoneByDoctor
// @Summary create milestone
// @Tags Process
// @Description the doctor create milestone (type: array)
// @Param RegistrationID body uint true "registration's ID"
// @Param Activity body string true "milestone's activity"
// @Produce json
// @Success 204 {string} api.ReturnedData{}
// @Router /milestone [POST]
func (h *ProcessHandler) CreateMileStoneByDoctor(c echo.Context) error {
	DB := utils.GetDB()

	RegistrationID, err := strconv.ParseUint(c.QueryParam("RegistrationID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}
	milestone := MileStone{
		RegistrationID: uint(RegistrationID),
		Activity:       c.QueryParam("Activity"),
	}

	result := DB.Create(&milestone)
	if result.Error != nil {
		return c.JSON(http.StatusUnprocessableEntity, api.Return("error", result.Error))
	}

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
func (h *ProcessHandler) UpdateMileStoneByDoctor(c echo.Context) error {
	DB := utils.GetDB()

	var Milestone MileStone
	var err error
	var Checked bool

	DB.First(&Milestone, c.QueryParam("MileStoneID"))
	Checked, err = strconv.ParseBool(c.Param("Checked"))

	if err != nil {
		log.Error("UpdateMileStoneByDoctor failed due to", err)
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}
	Milestone.Checked = Checked
	Milestone.Activity = c.QueryParam("Activity")

	return c.JSON(http.StatusOK, api.Return("ok", nil))
}

// DeleteMileStoneByDoctor
// @Summary delete milestone
// @Tags Process
// @Description the doctor delete milestone
// @Param MileStoneID path uint true "milestone's ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone/{MileStoneID} [DELETE]
func (h *ProcessHandler) DeleteMileStoneByDoctor(c echo.Context) error {
	// TODO verify identity
	DB := utils.GetDB()
	DB.Delete(&MileStone{}, c.Param("MileStoneID"))

	return c.JSON(http.StatusOK, api.Return("ok", nil))
}

// NOTE: use delete method, because there is no dependencies on MileStone
