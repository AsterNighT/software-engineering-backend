package cases

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
)

type CaseHandler struct {
}

// @Summary New a case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body string true "patient ID, doctor ID and department name"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case [POST]
func (h *CaseHandler) NewCase(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("NewCase", nil))
}

// @Summary Get case by a case-id list
// @Description
// @Tags Case
// @Produce json
// @Param  id query []uint true "case IDs"
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /case/{caseid} [GET]
func (h *CaseHandler) GetCaseInOneList(c echo.Context) error {
	// what is in here does not matter for now, just filling the swaggo comments and name the function
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetCaseInOneList", nil))
}

// @Summary Get prevoius cases by current case ID
// @Description
// @Tags Case
// @Produce json
// @Param  id path uint true "current case ID"
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /case/{currentcaseid} [GET]
func (h *CaseHandler) GetPreviousCaseByID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPreviousCaseByID", nil))
}

// @Summary Get lastest case by department name
// @Description
// @Tags Case
// @Produce json
// @Param  department path string true "department name"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/department [GET]
func (h *CaseHandler) GetPreviousCaseByDepartment(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPreviousCaseByDepartment", nil))
}

// @Summary Update case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body string true "case ID and updated details"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/{caseid} [PUT]
func (h *CaseHandler) UpdateCase(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("UpdateCase", nil))
}

// @Summary New a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionDetail body string true "case ID and prescription details"
// @Success 200 {object} api.ReturnedData{data=Prescription}
// @Router /case/{caseid}/prescription [POST]
func (h *CaseHandler) NewPrescription(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("NewPrescription", nil))
}

// @Summary Get prescription by prescription id
// @Description
// @Tags Case
// @Produce json
// @Param id path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{data=Prescription}
// @Router /case/prescription/{prescriptionid}  [GET]
func (h *CaseHandler) GetPrescriptionByPrescriptionID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPrescriptionByPrescriptionID", nil))
}

// @Summary Get prescriptions by case id
// @Description
// @Tags Case
// @Produce json
// @Param id path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=[]Prescription}
// @Router /case/{caseid}/prescription  [GET]
func (h *CaseHandler) GetPrescriptionByCaseID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPrescriptionByCaseID", nil))
}
