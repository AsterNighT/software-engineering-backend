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
// @Param pid formData uint true "patient ID"
// @Param did formData uint true "doctor ID"
// @Param department formData string true "department name"
// @Success 200 {string} string "{"msg":"establishment success"}"
// @Router /case [POST]
func (h *CaseHandler) NewCase(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, "establishment success")
}

// @Summary Get case by case id
// @Description
// @Tags Case
// @Produce json
// @Param  id path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/{id} [GET]
func (h *CaseHandler) GetCaseByCaseID(c echo.Context) error {
	// what is in here does not matter for now, just filling the swaggo comments and name the function
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetCaseByCaseID", nil))
}

// @Summary Get prevoius case by current case ID
// @Description
// @Tags Case
// @Produce json
// @Param  id path uint true "current case ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/{cid} [GET]
func (h *CaseHandler) GetPreviousCaseByID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPreviousCaseByID", nil))
}

// @Summary Get lastest case by department name
// @Description
// @Tags Case
// @Produce json
// @Param  department path uint true "department name"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /case/department [GET]
func (h *CaseHandler) GetPreviousCaseByDepartment(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPreviousCaseByDepartment", nil))
}

// @Summary Update complaint
// @Description
// @Tags Case
// @Produce json
// @Param id formData uint true "case ID"
// @Param details formData string true "complaint details"
// @Success 200 {string} string "{"msg":"modification success"}"
// @Router /case/complaint/{id} [PATCH]
func (h *CaseHandler) UpdateComplaint(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, "modification success")
}

// @Summary Update diagnosis
// @Description
// @Tags Case
// @Produce json
// @Param id formData uint true "case ID"
// @Param details formData string true "diagnosis details"
// @Success 200 {string} string "{"msg":"modification success"}"
// @Router /case/diagnosis/{id} [PATCH]
func (h *CaseHandler) UpdateDiagnosis(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, "modification success")
}

// @Summary Update Past History
// @Description
// @Tags Case
// @Produce json
// @Param id formData uint true "case ID"
// @Param details formData string true "Past History details"
// @Success 200 {string} string "{"msg":"modification success"}"
// @Router /case/PastHistory/{id} [PATCH]
func (h *CaseHandler) UpdatePastHistory(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, "modification success")
}

// @Summary New a prescrition
// @Description
// @Tags Prescription
// @Produce json
// @Param id formData uint true "case ID"
// @Param details formData string true "prescription details"
// @Success 200 {string} string "{"msg":"establishment success"}"
// @Router /prescription [POST]
func (h *CaseHandler) NewPrescription(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, "establishment success")
}

// @Summary Get prescription by prescription id
// @Description
// @Tags Prescription
// @Produce json
// @Param id path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{data=Prescription}
// @Router /prescription/{id}  [GET]
func (h *CaseHandler) GetPrescriptionByPrescriptionID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPrescriptionByPrescriptionID", nil))
}
