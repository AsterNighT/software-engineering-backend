package process

import (
	"encoding/json"
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/database/models"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"github.com/deckarep/golang-set"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

type ProcessHandler struct{}


func (h *ProcessHandler) Search(c echo.Context) error {
	response, err := http.Get(string("http://zhouxunwang.cn/data/?id=111&key=" + models.KevinKey + "&title=") + c.Param("KeyWord"))
	if err != nil{
		c.Logger().Debug("search failed...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.SearchFailed))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Debug("search failed...")
		}
	}(response.Body) //在回复后必须关闭回复的主体
	set :=  mapset.NewSet()
	resMap := Transformation(response)
	if v, ok := resMap["result"]; ok {
		list := v.([]interface{})
		for _, v := range list {
			nodeMap := v.(map[string]interface{})
			if verStr, ok := nodeMap["jzks"]; ok {
				strArray := strings.Fields(verStr.(string))
				for i := range strArray {
					if !set.Contains(strArray[i]) {
						set.Add(strArray[i])
					}
				}
			}
		}
	}
	names := make([]string, 1)
	it := set.Iterator()
	for elem := range it.C {
		names = append(names, elem.(string))
	}
	//fmt.Println(names)
	db := utils.GetDB()
	var departments []models.Department
	db.Where("name in (?)", names).Find(&departments)
	return c.JSON(http.StatusOK, api.Return("ok", departments))
}


// GetAllDepartments
// @Summary get all departments
// @Tags Process
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]models.Department}
// @Router /departments [GET]
func (h *ProcessHandler) GetAllDepartments(c echo.Context) error {
	db := utils.GetDB()

	// get all departments
	var departments []models.Department
	db.Find(&departments)

	return c.JSON(http.StatusOK, api.Return("ok", departments))
}

// GetDepartmentByID
// @Summary get a department by its ID
// @Tags Process
// @Description return a department's details by its ID
// @Param departmentID path uint true "department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=models.DepartmentDetailJSON}
// @Router /department/{DepartmentID} [GET]
func (h *ProcessHandler) GetDepartmentByID(c echo.Context) error {
	db := utils.GetDB()
	var department models.Department          // department basic information
	var schedules []models.DepartmentSchedule // schedules
	var doctorsAll = make([]models.Doctor, 0) // doctor total information
	var doctors = make([]string, 0)           // doctor's name (show to frontend)

	err := db.First(&department, c.Param("departmentID")).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Return("error", models.DepartmentNotFound))
	}
	db.Where("department_id = ?", department.ID).Find(&schedules)
	db.Where("department_id = ?", department.ID).Find(&doctorsAll)

	for _, doctor := range doctorsAll {
		var a models.Account
		err = db.First(&a, doctor.AccountID).Error
		if err == nil {
			doctors = append(doctors, a.LastName+a.FirstName)
		}
	}

	// wrap the return data
	returnDepartment := models.DepartmentDetailJSON{
		ID:        department.ID,
		Name:      department.Name,
		Detail:    department.Detail,
		Doctors:   doctors,
		Schedules: schedules,
	}

	c.Logger().Debug("GetDepartmentByID")
	return c.JSON(http.StatusOK, api.Return("ok", returnDepartment))
}

// CreateRegistrationTX
// @Summary create registration
// @Tags Process
// @Description return registration state
// @Param department_id body uint true "department ID"
// @Param year body int true "Year"
// @Param month body int true "Month"
// @Param day body int true "Day"
// @Param halfday body int true "HalfDay"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=int}
// @Router /registrations [POST]
func (h *ProcessHandler) CreateRegistrationTX(c echo.Context) error {
	type RegistrationSubmitJSON struct {
		DepartmentID uint               `json:"department_id"`
		Year         int                `json:"year"`
		Month        int                `json:"month"`
		Day          int                `json:"day"`
		HalfDay      models.HalfDayEnum `json:"halfday" validate:"halfday"`
	}

	// extract submit data
	var submit RegistrationSubmitJSON
	if err := c.Bind(&submit); err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.InvalidSubmitFormat))
	}

	// validate halfday enum
	err := models.Validate.Struct(submit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.InvalidSubmitFormat))
	}

	var res = 1

	db := utils.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		// get department
		var department models.Department
		err = db.First(&department, submit.DepartmentID).Error
		if err != nil {
			c.Logger().Error(models.DepartmentNotFound)
			return err
		}

		// get patient
		var patient models.Patient
		err = db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			c.Logger().Error(models.PatientNotFound)
			return err
		}

		var possibleDuplicates []models.Registration

		//check duplicate registration
		if db.Where(&models.Registration{
			DepartmentID: department.ID,
			PatientID:    patient.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}).Find(&possibleDuplicates).RowsAffected > 0 {
			for i := range possibleDuplicates {
				if possibleDuplicates[i].Status != models.Terminated {
					return echo.ErrBadRequest
				}
			}
		}

		// check schedule
		var schedule models.DepartmentSchedule
		err = db.Where(&models.DepartmentSchedule{
			DepartmentID: department.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}).First(&schedule).Error

		// invalid schedule return and unlock
		if err != nil || !models.ValidateSchedule(&schedule) {
			return err
		} else if schedule.Current >= schedule.Capacity {
			c.Logger().Error(models.NotEnoughCapacity)
			return err
		}

		// assign the doctor with the minimal registrations
		var doctors []models.Doctor
		db.Where("department_id = ?", department.ID).Find(&doctors)
		var doctorRegistrationCount = make([]int64, len(doctors))

		registration := models.Registration{
			DepartmentID: department.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}

		// find min count of registrations
		minCount, minIndex := int64(math.MaxInt64), -1
		for i := range doctors {
			db.Model(&models.Registration{}).Where(&registration).Count(&doctorRegistrationCount[i])
			if minCount > doctorRegistrationCount[i] {
				minCount = doctorRegistrationCount[i]
				minIndex = i
			}
		}

		// cannot find a doctor
		if minIndex == -1 {
			c.Logger().Error(models.CannotAssignDoctor)
			return echo.ErrBadRequest
		}

		// process schedule
		registration.PatientID = patient.ID
		registration.DoctorID = doctors[minIndex].ID

		if err := db.Create(&registration).Error; err != nil {
			return err
		}

		// create corresponding case
		var correspondingCase models.Case
		var preCase models.Case

		correspondingCase.Registration = registration
		correspondingCase.DoctorID = doctors[minIndex].ID
		correspondingCase.PatientID = patient.ID
		correspondingCase.Department = department.Name

		dbErr := db.Where("department = ?", department.Name).Order("date DESC").Limit(1).First(&preCase)
		if dbErr == nil {
			correspondingCase.PreviousCase = &preCase
		}

		db.Create(&correspondingCase)

		schedule.Current = schedule.Current + 1
		if err := db.Save(&schedule).Error; err != nil {
			c.Logger().Error(models.InvalidSchedule)
			return err
		}

		// get name
		var patientAccount models.Account
		var doctor models.Doctor
		var doctorAccount models.Account
		db.First(&patientAccount, registration.PatientID)
		db.First(&doctor, registration.DoctorID)
		db.First(&doctorAccount, doctor.AccountID)

		res = int(registration.ID)

		return nil
	})

	if err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.CreateRegistrationFailed))
	}


	return c.JSON(http.StatusOK, api.Return("ok", res))
}

// GetRegistrations
// @Summary get all registrations
// @Tags Process
// @Description display all registrations
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]models.RegistrationJSON}
// @Router /registrations [GET]
func (h *ProcessHandler) GetRegistrations(c echo.Context) error {
	db := utils.GetDB()
	var registrations []models.Registration

	// get account type
	var acc models.Account
	err := db.First(&acc, c.Get("id")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	if acc.Type == models.PatientType {
		// get patient
		var patient models.Patient
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.PatientNotFound))
		}

		db.Where("patient_id = ?", patient.ID).Find(&registrations)
	} else if acc.Type == models.DoctorType {
		// get doctor
		var doctor models.Doctor
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.DoctorNotFound))
		}

		db.Where("doctor_id = ?", doctor.ID).Find(&registrations)
	} else {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	var registrationJSONs = make([]models.RegistrationJSON, len(registrations))

	for i := range registrations {
		var department models.Department
		db.First(&department, registrations[i].DepartmentID)
		registrationJSONs[i] = models.RegistrationJSON{
			ID:         registrations[i].ID,
			Department: department.Name,
			Status:     registrations[i].Status,
			Year:       registrations[i].Year,
			Month:      registrations[i].Month,
			Day:        registrations[i].Day,
			HalfDay:    registrations[i].HalfDay,
		}
	}
	c.Logger().Debug("GetRegistrationsByPatient")
	return c.JSON(http.StatusOK, api.Return("ok", registrationJSONs))
}

// GetRegistrationByID
// @Summary get a registration by its ID
// @Tags Process
// @Description return a registration details by its ID
// @Param registrationID path uint true "registration's ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=models.Registration}
// @Router /registration/{registrationID} [GET]
func (h *ProcessHandler) GetRegistrationByID(c echo.Context) error {
	db := utils.GetDB()
	var registration models.Registration

	// get account type
	var acc models.Account
	err := db.First(&acc, c.Get("id")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	// judge account type
	if acc.Type == models.PatientType {
		// get patient
		var patient models.Patient
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.PatientNotFound))
		}

		err = db.Where("patient_id = ?", patient.ID).First(&registration, c.Param("registrationID")).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.RegistrationNotFound))
		}
	} else if acc.Type == models.DoctorType {
		// get doctor
		var doctor models.Doctor
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.DoctorNotFound))
		}

		err = db.Where("doctor_id = ?", doctor.ID).First(&registration, c.Param("registrationID")).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", models.RegistrationNotFound))
		}
	} else {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	// get names
	var department models.Department
	var patientAccount models.Account
	var doctor models.Doctor
	var doctorAccount models.Account
	db.First(&department, registration.DepartmentID)
	db.First(&patientAccount, registration.PatientID)
	db.First(&doctor, registration.DoctorID)
	db.First(&doctorAccount, doctor.AccountID)

	registrationJSON := models.RegistrationDetailJSON{
		ID:              registration.ID,
		Department:      department.Name,
		Patient:         patientAccount.LastName + patientAccount.FirstName,
		Doctor:          doctorAccount.LastName + doctorAccount.FirstName,
		Year:            registration.Year,
		Month:           registration.Month,
		Day:             registration.Day,
		HalfDay:         registration.HalfDay,
		Status:          registration.Status,
		TerminatedCause: registration.TerminatedCause,
	}

	// get milestones
	var milestones []models.MileStone
	db.Where("registration_id = ?", registration.ID).Find(&milestones)
	registrationJSON.MileStone = milestones

	c.Logger().Debug("GetRegistrationByDoctor")
	return c.JSON(http.StatusOK, api.Return("ok", registrationJSON))
}

// UpdateRegistrationStatus
// @Summary update registration status
// @Tags Process
// @Description update registration status
// @Param registrationID path uint true "registration ID"
// @Param status body string true "next status of current registration"
// @Param terminatedCause body string false "if registration is shutdown, a cause is required"
// @Produce json
// @Success 200 {object} api.ReturnedData{}
// @Router /registration/{registrationID} [PUT]
func (h *ProcessHandler) UpdateRegistrationStatus(c echo.Context) error {

	type RegistrationSubmitJSON struct {
		RegistrationStatus string `json:"status"`
		TerminatedCause    string `json:"terminatedCause"`
	}

	// extract submit data
	var submit RegistrationSubmitJSON
	if err := c.Bind(&submit); err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.InvalidSubmitFormat))
	}

	db := utils.GetDB()
	// validation for status: select in frontend
	status := models.RegistrationStatusEnum(submit.RegistrationStatus)
	terminatedCause := submit.TerminatedCause
	var registration models.Registration
	err := db.First(&registration, c.Param("registrationID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}
	currentStatus := registration.Status
	registration.Status = status
	var acc models.Account
	err = db.Where("id = ?", c.Get("id").(uint)).First(&acc).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.AccountNotFound))
	}
	if currentStatus == models.Committed {
		if acc.Type == models.PatientType {
			if status == models.Terminated {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
		}
		if acc.Type == models.DoctorType {
			if status == models.Accepted {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
			if status == models.Terminated {
				if terminatedCause != "" {
					registration.Status = status
					registration.TerminatedCause = terminatedCause
					db.Save(&registration)
					return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
				}
				return c.JSON(http.StatusBadRequest, api.Return("failed", "Missing terminated causes"))
			}
		}
	}
	if currentStatus == models.Accepted {
		if acc.Type == "doctor" {
			if status == models.Terminated {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
		}
	}
	return c.JSON(http.StatusBadRequest, api.Return("failed", models.RegistrationUpdateFailed))
}

// CreateMileStoneByDoctor
// @Summary create milestone
// @Tags Process
// @Description the doctor create milestone (type: array)
// @Param registration_id body uint true "registration's ID"
// @Param activity body string true "milestone's activity"
// @Produce json
// @Success 204 {string} api.ReturnedData{}
// @Router /milestones [POST]
func (h *ProcessHandler) CreateMileStoneByDoctor(c echo.Context) error {
	type MileStoneSubmitJSON struct {
		RegistrationID uint   `json:"registration_id"`
		Activity       string `json:"activity"`
	}

	// extract submit data
	var submit MileStoneSubmitJSON
	if err := c.Bind(&submit); err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.InvalidSubmitFormat))
	}

	db := utils.GetDB()

	// get doctor
	var doctor models.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.DoctorNotFound))
	}

	var registration models.Registration
	err = db.First(&registration, submit.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.RegistrationNotFound))
	} else if registration.Status == models.Terminated {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.MileStoneUnauthorized))
	}

	mileStone := models.MileStone{
		RegistrationID: submit.RegistrationID,
		Activity:       submit.Activity,
	}

	err = db.Create(&mileStone).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.CreateMileStoneFailed))
	}

	return c.JSON(http.StatusOK, api.Return("ok", nil))
}

// UpdateMileStoneByDoctor
// @Summary update milestone
// @Tags Process
// @Description the doctor update milestone (check milestone)
// @Param mileStoneID path uint true "milestone's ID"
// @Param activity body string false "updated milestone's activity"
// @Param checked body boolean true "milestone is checked or not"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone/{mileStoneID} [PUT]
func (h *ProcessHandler) UpdateMileStoneByDoctor(c echo.Context) error {
	type MileStoneSubmitJSON struct {
		Activity string `json:"activity"`
		Checked  bool   `json:"checked"`
	}

	// extract submit data
	var submit MileStoneSubmitJSON
	if err := c.Bind(&submit); err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", models.InvalidSubmitFormat))
	}

	db := utils.GetDB()
	// get doctor
	var doctor models.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.DoctorNotFound))
	}

	// get mileStone
	var mileStone models.MileStone
	err = db.First(&mileStone, c.Param("mileStoneID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.MileStoneNotFound))
	}
	var registration models.Registration
	err = db.First(&registration, mileStone.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.RegistrationNotFound))
	}

	// check milestone authority
	if registration.Status == models.Terminated || registration.DoctorID != doctor.ID {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.MileStoneUnauthorized))
	}

	mileStone.Activity = submit.Activity
	mileStone.Checked = submit.Checked

	db.Save(&mileStone)
	return c.JSON(http.StatusOK, api.Return("ok", nil))

}

// DeleteMileStoneByDoctor
// @Summary delete milestone
// @Tags Process
// @Description the doctor delete milestone
// @Param mileStoneID path uint true "milestone's ID"
// @Produce json
// @Success 200 {string} api.ReturnedData{}
// @Router /milestone/{mileStoneID} [DELETE]
func (h *ProcessHandler) DeleteMileStoneByDoctor(c echo.Context) error {
	db := utils.GetDB()

	// get doctor
	var doctor models.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.DoctorNotFound))
	}

	// get mileStone
	var mileStone models.MileStone
	err = db.First(&mileStone, c.Param("mileStoneID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.MileStoneNotFound))
	}
	var registration models.Registration
	err = db.First(&registration, mileStone.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.RegistrationNotFound))
	}

	// check milestone authority
	if registration.Status == models.Terminated || registration.DoctorID != doctor.ID {
		return c.JSON(http.StatusBadRequest, api.Return("error", models.MileStoneUnauthorized))
	}

	db.Delete(&mileStone)
	return c.JSON(http.StatusOK, api.Return("ok", nil))
}

func Transformation(response *http.Response) map[string]interface{}{
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil
		}
	}
	return result
}
