package process

import (
	"fmt"
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	//"github.com/deckarep/golang-set"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"math"
	"net/http"
)

type ProcessHandler struct{}


func (h *ProcessHandler) Search(c echo.Context) error {
	response, err := http.Get(string("http://zhouxunwang.cn/data/?id=111&key=" + KevinKey + "&title=") + c.Param("KeyWord"))
	if err != nil{
		c.Logger().Debug("search failed...")
		return c.JSON(http.StatusBadRequest, api.Return("error", SearchFailed))
	}
	defer response.Body.Close()//在回复后必须关闭回复的主体

	resMap := Transformation(response)
	if v, ok := resMap["result"]; ok {
		list := v.([]interface{})
		for _, v := range list {
			nodeMap := v.(map[string]interface{})
			if verStr, ok := nodeMap["jzks"]; ok {
				fmt.Println(verStr)
			}
		}
	}

	//db := utils.GetDB()
	//var departments []Department
	//db.Where("department = ?", ).Find(&departments)
	//return c.JSON(http.StatusOK, api.Return("ok", departments))
	return nil
}


// GetAllDepartments
// @Summary get all departments
// @Tags Process
// @Description display all departments of a hospital
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]Department}
// @Router /departments [GET]
func (h *ProcessHandler) GetAllDepartments(c echo.Context) error {
	db := utils.GetDB()

	// get all departments
	var departments []Department
	db.Find(&departments)

	return c.JSON(http.StatusOK, api.Return("ok", departments))
}

// GetDepartmentByID
// @Summary get a department by its ID
// @Tags Process
// @Description return a department's details by its ID
// @Param departmentID path uint true "department ID"
// @Produce json
// @Success 200 {object} api.ReturnedData{data=DepartmentDetailJSON}
// @Router /department/{DepartmentID} [GET]
func (h *ProcessHandler) GetDepartmentByID(c echo.Context) error {
	db := utils.GetDB()
	var department Department                  // department basic information
	var schedules []DepartmentSchedule         // schedules
	var doctorsAll = make([]account.Doctor, 0) // doctor total information
	var doctors = make([]string, 0)            // doctor's name (show to frontend)

	err := db.First(&department, c.Param("departmentID")).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Return("error", DepartmentNotFound))
	}
	db.Where("department_id = ?", department.ID).Find(&schedules)
	db.Where("department_id = ?", department.ID).Find(&doctorsAll)

	for _, doctor := range doctorsAll {
		var a account.Account
		err = db.First(&a, doctor.AccountID).Error
		if err == nil {
			doctors = append(doctors, a.LastName+a.FirstName)
		}
	}

	// wrap the return data
	returnDepartment := DepartmentDetailJSON{
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
		DepartmentID uint        `json:"department_id"`
		Year         int         `json:"year"`
		Month        int         `json:"month"`
		Day          int         `json:"day"`
		HalfDay      HalfDayEnum `json:"halfday" validate:"halfday"`
	}

	// extract submit data
	var submit RegistrationSubmitJSON
	if err := c.Bind(&submit); err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", InvalidSubmitFormat))
	}

	// validate halfday enum
	err := validate.Struct(submit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", InvalidSubmitFormat))
	}

	var res = 1

	db := utils.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		// get department
		var department Department
		err = db.First(&department, submit.DepartmentID).Error
		if err != nil {
			c.Logger().Error(DepartmentNotFound)
			return err
		}

		// get patient
		var patient account.Patient
		err = db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			c.Logger().Error(PatientNotFound)
			return err
		}

		var possibleDuplicates []Registration

		//check duplicate registration
		if db.Where(&Registration{
			DepartmentID: department.ID,
			PatientID:    patient.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}).Find(&possibleDuplicates).RowsAffected > 0 {
			for i := range possibleDuplicates {
				if possibleDuplicates[i].Status != terminated {
					return echo.ErrBadRequest
				}
			}
		}

		// check schedule
		var schedule DepartmentSchedule
		err = db.Where(&DepartmentSchedule{
			DepartmentID: department.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}).First(&schedule).Error

		// invalid schedule return and unlock
		if err != nil || !validateSchedule(&schedule) {
			return err
		} else if schedule.Current >= schedule.Capacity {
			c.Logger().Error(NotEnoughCapacity)
			return err
		}

		// assign the doctor with the minimal registrations
		var doctors []account.Doctor
		db.Where("department_id = ?", department.ID).Find(&doctors)
		var doctorRegistrationCount = make([]int64, len(doctors))

		registration := Registration{
			DepartmentID: department.ID,
			Year:         submit.Year,
			Month:        submit.Month,
			Day:          submit.Day,
			HalfDay:      submit.HalfDay,
		}

		// find min count of registrations
		minCount, minIndex := int64(math.MaxInt64), -1
		for i := range doctors {
			db.Model(&Registration{}).Where(&registration).Count(&doctorRegistrationCount[i])
			if minCount > doctorRegistrationCount[i] {
				minCount = doctorRegistrationCount[i]
				minIndex = i
			}
		}

		// cannot find a doctor
		if minIndex == -1 {
			c.Logger().Error(CannotAssignDoctor)
			return echo.ErrBadRequest
		}

		// process schedule
		registration.PatientID = patient.ID
		registration.DoctorID = doctors[minIndex].ID

		if err := db.Create(&registration).Error; err != nil {
			return err
		}

		schedule.Current = schedule.Current + 1
		if err := db.Save(&schedule).Error; err != nil {
			c.Logger().Error(InvalidSchedule)
			return err
		}

		// get name
		var patientAccount account.Account
		var doctor account.Doctor
		var doctorAccount account.Account
		db.First(&patientAccount, registration.PatientID)
		db.First(&doctor, registration.DoctorID)
		db.First(&doctorAccount, doctor.AccountID)

		res = int(registration.ID)

		return nil
	})

	if err != nil {
		c.Logger().Debug("JSON format failed when trying to create a registration ...")
		return c.JSON(http.StatusBadRequest, api.Return("error", CreateRegistrationFailed))
	}


	return c.JSON(http.StatusOK, api.Return("ok", res))
}

// GetRegistrations
// @Summary get all registrations
// @Tags Process
// @Description display all registrations
// @Produce json
// @Success 200 {object} api.ReturnedData{data=[]RegistrationJSON}
// @Router /registrations [GET]
func (h *ProcessHandler) GetRegistrations(c echo.Context) error {
	db := utils.GetDB()
	var registrations []Registration

	// get account type
	var acc account.Account
	err := db.First(&acc, c.Get("id")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	if acc.Type == account.PatientType {
		// get patient
		var patient account.Patient
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", PatientNotFound))
		}

		db.Where("patient_id = ?", patient.ID).Find(&registrations)
	} else if acc.Type == account.DoctorType {
		// get doctor
		var doctor account.Doctor
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", DoctorNotFound))
		}

		db.Where("doctor_id = ?", doctor.ID).Find(&registrations)
	} else {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	var registrationJSONs = make([]RegistrationJSON, len(registrations))

	for i := range registrations {
		var department Department
		db.First(&department, registrations[i].DepartmentID)
		registrationJSONs[i] = RegistrationJSON{
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
// @Success 200 {object} api.ReturnedData{data=Registration}
// @Router /registration/{registrationID} [GET]
func (h *ProcessHandler) GetRegistrationByID(c echo.Context) error {
	db := utils.GetDB()
	var registration Registration

	// get account type
	var acc account.Account
	err := db.First(&acc, c.Get("id")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	// judge account type
	if acc.Type == account.PatientType {
		// get patient
		var patient account.Patient
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&patient).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", PatientNotFound))
		}

		err = db.Where("patient_id = ?", patient.ID).First(&registration, c.Param("registrationID")).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", RegistrationNotFound))
		}
	} else if acc.Type == account.DoctorType {
		// get doctor
		var doctor account.Doctor
		err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", DoctorNotFound))
		}

		err = db.Where("doctor_id = ?", doctor.ID).First(&registration, c.Param("registrationID")).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Return("error", RegistrationNotFound))
		}
	} else {
		return c.JSON(http.StatusUnauthorized, api.Return("error", nil))
	}

	// get names
	var department Department
	var patientAccount account.Account
	var doctor account.Doctor
	var doctorAccount account.Account
	db.First(&department, registration.DepartmentID)
	db.First(&patientAccount, registration.PatientID)
	db.First(&doctor, registration.DoctorID)
	db.First(&doctorAccount, doctor.AccountID)

	registrationJSON := RegistrationDetailJSON{
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
	var milestones []MileStone
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
		return c.JSON(http.StatusBadRequest, api.Return("error", InvalidSubmitFormat))
	}

	db := utils.GetDB()
	// validation for status: select in frontend
	status := RegistrationStatusEnum(submit.RegistrationStatus)
	terminatedCause := submit.TerminatedCause
	var registration Registration
	err := db.First(&registration, c.Param("registrationID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}
	currentStatus := registration.Status
	registration.Status = status
	var acc account.Account
	err = db.Where("id = ?", c.Get("id").(uint)).First(&acc).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", AccountNotFound))
	}
	if currentStatus == committed {
		if acc.Type == account.PatientType {
			if status == terminated {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
		}
		if acc.Type == account.DoctorType {
			if status == accepted {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
			if status == terminated {
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
	if currentStatus == accepted {
		if acc.Type == "doctor" {
			if status == terminated {
				registration.Status = status
				db.Save(&registration)
				return c.JSON(http.StatusOK, api.Return("ok", "修改挂号成功"))
			}
		}
	}
	return c.JSON(http.StatusBadRequest, api.Return("failed", RegistrationUpdateFailed))
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
		return c.JSON(http.StatusBadRequest, api.Return("error", InvalidSubmitFormat))
	}

	db := utils.GetDB()

	// get doctor
	var doctor account.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", DoctorNotFound))
	}

	var registration Registration
	err = db.First(&registration, submit.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", RegistrationNotFound))
	} else if registration.Status == terminated {
		return c.JSON(http.StatusBadRequest, api.Return("error", MileStoneUnauthorized))
	}

	mileStone := MileStone{
		RegistrationID: submit.RegistrationID,
		Activity:       submit.Activity,
	}

	err = db.Create(&mileStone).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", CreateMileStoneFailed))
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
		return c.JSON(http.StatusBadRequest, api.Return("error", InvalidSubmitFormat))
	}

	db := utils.GetDB()
	// get doctor
	var doctor account.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", DoctorNotFound))
	}

	// get mileStone
	var mileStone MileStone
	err = db.First(&mileStone, c.Param("mileStoneID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", MileStoneNotFound))
	}
	var registration Registration
	err = db.First(&registration, mileStone.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", RegistrationNotFound))
	}

	// check milestone authority
	if registration.Status == terminated || registration.DoctorID != doctor.ID {
		return c.JSON(http.StatusBadRequest, api.Return("error", MileStoneUnauthorized))
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
	var doctor account.Doctor
	err := db.Where("account_id = ?", c.Get("id").(uint)).First(&doctor).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", DoctorNotFound))
	}

	// get mileStone
	var mileStone MileStone
	err = db.First(&mileStone, c.Param("mileStoneID")).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", MileStoneNotFound))
	}
	var registration Registration
	err = db.First(&registration, mileStone.RegistrationID).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", RegistrationNotFound))
	}

	// check milestone authority
	if registration.Status == terminated || registration.DoctorID != doctor.ID {
		return c.JSON(http.StatusBadRequest, api.Return("error", MileStoneUnauthorized))
	}

	db.Delete(&mileStone)
	return c.JSON(http.StatusOK, api.Return("ok", nil))
}
