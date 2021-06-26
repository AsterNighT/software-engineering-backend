package cases

import (
	"strconv"
	"time"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CaseHandler struct {
}

type MedicineHandler struct {
}

// @Summary Get all cases
// @Description
// @Tags Case
// @Produce json
// @Param department query string false "department name" nil
// @Param patientID query uint false "patient ID" nil
// @Param doctorID query uint false "doctor ID" nil
// @Param before query string false "a timestamp marking end time without timezone" nil
// @Param after query string false "a timestamp marking start time without timezone" nil
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /cases [GET]
func (h *CaseHandler) GetAllCases(c echo.Context) error {

	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()

	if c.QueryParam(("patientID")) != "" {
		db = db.Where("patient_id = ?", c.QueryParam("patientID"))
	}
	if c.QueryParam(("doctorID")) != "" {
		db = db.Where("doctor_id = ?", c.QueryParam("doctorID"))
	}
	if c.QueryParam("department") != "" {
		db = db.Where("department LIKE ?", "%"+c.QueryParam("department")+"%")
	}
	if c.QueryParam("before") != "" {
		db = db.Where("date <= ?", c.QueryParam("before"))
	}
	if c.QueryParam("after") != "" {
		db = db.Where("date >= ?", c.QueryParam("after"))
	}

	var cases []Case
	db.Preload("Prescriptions").Preload("Prescriptions.Guidelines").Preload("Prescriptions.Guidelines.Medicine").Find(&cases)

	c.Logger().Debug("GetAllCases")
	return c.JSON(200, api.Return("ok", cases))
}

// @Summary Get case by caseID
// @Description
// @Tags Case
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param caseID path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /patient/{patientID}/cases/{caseID} [GET]
func (h *CaseHandler) GetCaseByCaseID(c echo.Context) error {

	var patID int
	patID, _ = strconv.Atoi(c.Param("patientID"))
	if !FromPatient(c, uint(patID)) {
		return c.JSON(403, api.Return("you cannot access this case", nil))
	}

	db := utils.GetDB()
	var case1 Case
	db.Preload("Prescriptions").Preload("Prescriptions.Guidelines").Preload("Prescriptions.Guidelines.Medicine").First(&case1, c.Param("caseID"))

	c.Logger().Debug("GetCasebyCaseID")
	return c.JSON(200, api.Return("ok", case1))
}

// @Summary Get the last case
// @Description
// @Tags Case
// @Produce json
// @Param patientID path uint true "patient ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /patient/{patientID}/case [GET]
func (h *CaseHandler) GetLastCaseByPatientID(c echo.Context) error {

	var patID int
	patID, _ = strconv.Atoi(c.Param("patientID"))
	if !FromPatient(c, uint(patID)) {
		return c.JSON(403, api.Return("you cannot access this case", nil))
	}

	db := utils.GetDB()
	db = db.Where("patient_id = ?", c.Param("patientID"))

	var case1 Case
	db.Preload("Prescriptions").Preload("Prescriptions.Guidelines").Preload("Prescriptions.Guidelines.Medicine").Order("date DESC").Limit(1).First(&case1)

	c.Logger().Debug("GetLastCase")
	return c.JSON(200, api.Return("ok", case1))
}

// @Summary Get cases by patient ID
// @Description
// @Tags Case
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param department query string false "department name" nil
// @Param doctorID query string false "doctor name" nil
// @Param before query uint false "a timestamp marking end time" nil
// @Param after query uint false "a timestamp marking start time" nil
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /patient/{patientID}/cases [GET]
func (h *CaseHandler) GetCasesByPatientID(c echo.Context) error {

	var patID int
	patID, _ = strconv.Atoi(c.Param("patientID"))
	if !FromPatient(c, uint(patID)) {
		return c.JSON(403, api.Return("you cannot access this case", nil))
	}

	db := utils.GetDB()
	db = db.Where("patient_id = ?", c.Param("patientID"))
	if c.QueryParam(("doctorID")) != "" {
		db = db.Where("doctor_id = ?", c.QueryParam("doctorID"))
	}
	if c.QueryParam("department") != "" {
		db = db.Where("department LIKE ?", "%"+c.QueryParam("department")+"%")
	}
	if c.QueryParam("before") != "" {
		db = db.Where("date <= ?", c.QueryParam("before"))
	}
	if c.QueryParam("after") != "" {
		db = db.Where("date >= ?", c.QueryParam("after"))
	}

	var cases []Case
	db.Preload("Prescriptions").Preload("Prescriptions.Guidelines").Preload("Prescriptions.Guidelines.Medicine").Find(&cases)

	c.Logger().Debug("GetCasesByPatientID")
	return c.JSON(200, api.Return("ok", cases))
}

// @Summary New a case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body Case true "patient ID, doctor ID, department name and other case details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case [POST]
func (h *CaseHandler) NewCase(c echo.Context) error {

	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()
	var cas Case
	err := utils.ExtractDataWithValidating(c, &cas)
	if err != nil {
		return c.JSON(400, api.Return("error", err.Error()))
	}
	var acc1, acc2 account.Account
	db.First(&acc1, cas.DoctorID)
	db.First(&acc2, cas.PatientID)
	cas.DoctorName = acc1.LastName + " " + acc1.FirstName
	cas.PatientName = acc2.LastName + " " + acc2.FirstName
	cas.Age = 18
	cas.Gender = "男"
	cas.Date = time.Now()
	result := db.Create(&cas)
	if result.Error != nil {
		return c.JSON(400, api.Return("error", result.Error.Error()))
	}
	c.Logger().Debug("NewCase")
	return c.JSON(200, api.Return("ok", cas.ID))
}

// @Summary Delete a case
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "case ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID} [DELETE]
func (h *CaseHandler) DeleteCaseByCaseID(c echo.Context) error {

	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()
	db.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Prescriptions.Guidelines.Medicine").Delete(&Case{}, c.Param("caseID"))
	c.Logger().Debug("DeleteCaseByCaseID")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get previous cases
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "current case ID"
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /patient/{patientID}/case/{caseID} [GET]
func (h *CaseHandler) GetPreviousCases(c echo.Context) error {
	db := utils.GetDB()
	var case1 Case
	var cases []Case
	db.First(&case1, c.Param("caseID"))
	if !FromPatient(c, case1.PatientID) {
		return c.JSON(403, api.Return("you cannot access this case", nil))
	}
	for case1.PreviousCaseID != nil {
		case1.ID = *case1.PreviousCaseID
		db.First(&case1, *case1.PreviousCaseID)
		cases = append(cases, case1)
	}
	c.Logger().Debug("GetPreviousCases")
	return c.JSON(200, api.Return("ok", cases))
}

// @Summary Update case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body Case true "case ID and updated details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID} [PUT]
func (h *CaseHandler) UpdateCase(c echo.Context) error {
	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()
	var cas Case
	err := utils.ExtractDataWithValidating(c, &cas)
	if err != nil {
		return c.JSON(400, api.Return("error", err.Error()))
	}
	result := db.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Prescriptions.Guidelines.Medicine").Model(&cas).Updates(cas)
	if result.Error != nil {
		return c.JSON(400, api.Return("error", result.Error.Error()))
	}
	c.Logger().Debug("UpdateCase")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary New a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionDetail body Prescription true "case ID and prescription details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID}/prescription [POST]
func (h *CaseHandler) NewPrescription(c echo.Context) error {
	db := utils.GetDB()
	var pre Prescription
	err := utils.ExtractDataWithValidating(c, &pre)
	if err != nil {
		return c.JSON(400, api.Return("error", err.Error()))
	}
	var oldCase Case
	db.First(&oldCase, pre.CaseID)
	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}
	result := db.Create(&pre)
	if result.Error != nil {
		return c.JSON(400, api.Return("error", result.Error.Error()))
	}
	c.Logger().Debug("NewPrescription")
	return c.JSON(200, api.Return("ok", pre.ID))
}

// @Summary Delete a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionID path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID} [DELETE]
func (h *CaseHandler) DeletePrescription(c echo.Context) error {

	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()
	db.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Guidelines.Medicine").Delete(&Prescription{}, c.Param("prescriptionID"))
	c.Logger().Debug("DeletePrescription")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Update a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionDetails body Prescription true "prescription ID and updated details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID} [PUT]
func (h *CaseHandler) UpdatePrescription(c echo.Context) error {

	if !FromDoctor(c) {
		return c.JSON(403, api.Return("only doctor can access this endpoint", nil))
	}

	db := utils.GetDB()
	var pre Prescription
	err := utils.ExtractDataWithValidating(c, &pre)
	if err != nil {
		return c.JSON(400, api.Return("error", err.Error()))
	}
	result := db.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Guidelines.Medicine").Model(&pre).Updates(pre)
	if result.Error != nil {
		return c.JSON(400, api.Return("error", result.Error.Error()))
	}
	c.Logger().Debug("UpdatePrescription")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get prescription by prescription id
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionID path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{data=Prescription}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID}  [GET]
func (h *CaseHandler) GetPrescriptionByPrescriptionID(c echo.Context) error {
	db := utils.GetDB()
	var pre Prescription
	db.Preload("Guidelines").Preload("Guidelines.Medicine").First(&pre, c.Param("prescriptionID"))
	var oldCase Case
	db.First(&oldCase, pre.CaseID)
	if !FromPatient(c, oldCase.PatientID) {
		return c.JSON(403, api.Return("you cannot access this prescription", nil))
	}
	// var jsonStu []byte
	// jsonStu, _ = json.Marshal(pre)
	// fmt.Println(string(jsonStu))

	c.Logger().Debug("GetPrescriptionByPrescriptionID")
	return c.JSON(200, api.Return("ok", pre))
}

// @Summary Get prescriptions by case id
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=[]Prescription}
// @Router /patient/{patientID}/case/{caseID}/prescription  [GET]
func (h *CaseHandler) GetPrescriptionByCaseID(c echo.Context) error {
	db := utils.GetDB()
	var pres []Prescription
	var oldCase Case
	db.First(&oldCase, c.Param("caseID"))
	if !FromPatient(c, oldCase.PatientID) {
		return c.JSON(403, api.Return("you cannot access this prescription", nil))
	}
	db.Where("case_id = ?", c.Param("caseID")).Preload("Guidelines").Preload("Guidelines.Medicine").Find(&pres)
	c.Logger().Debug("GetPrescriptionByCaseID")
	return c.JSON(200, api.Return("ok", pres))
}

// @Summary Query medicine by key word
// @Description
// @Tags Medicine
// @Produce json
// @Param q query string true "full-text search"
// @Success 200 {object} api.ReturnedData{data=[]Medicine}
// @Router /medicine [GET]
func (h *MedicineHandler) GetMedicines(c echo.Context) error {
	db := utils.GetDB()
	var meds []Medicine
	db.Where("name LIKE ?", "%"+c.QueryParam("q")+"%").Find(&meds)
	c.Logger().Debug("QueryMedicine")
	return c.JSON(200, api.Return("ok", meds))
}
