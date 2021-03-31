package cases

import (
	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
)

type CaseHandler struct {
}

// @Summary Get the lastest case
// @Description
// @Tags Case
// @Produce json
// @Param patientID body Case true "patient ID"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /patient/{patientID}/case [GET]
func (h *CaseHandler) GetLastCase(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetLastCase", nil))
}

// @Summary Get a case ID list by patient ID
// @Description
// @Tags Case
// @Produce json
// @Param patientID path uint true "patient ID"
// @Param department query string false "department name" nil
// @Success 200 {array} api.ReturnedData{data=[]uint}
// @Router /patient/{patientID} [GET]
func (h *CaseHandler) GetCasesByPatientID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetCasesByPatientID", nil))
}

// @Summary New a case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body Case true "patient ID, doctor ID, department name and other case details"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /patient/{patientID}/case [POST]
func (h *CaseHandler) NewCase(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("NewCase", nil))
}

// @Summary Delete a case
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "case ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID} [DELETE]
func (h *CaseHandler) DeleteCaseByCaseID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("DeleteCaseByCaseID", nil))
}

// @Summary Get prevoius cases
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "current case ID"
// @Success 200 {object} api.ReturnedData{data=[]Case}
// @Router /patient/{patientID}/case/{caseID} [GET]
func (h *CaseHandler) GetPreviousCases(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPreviousCases", nil))
}

// @Summary Update case
// @Description
// @Tags Case
// @Produce json
// @Param caseDetail body string true "case ID and updated details"
// @Success 200 {object} api.ReturnedData{data=Case}
// @Router /patient/{patientID}/case/{caseID} [PUT]
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
// @Router /patient/{patientID}/case/{caseID}/prescription [POST]
func (h *CaseHandler) NewPrescription(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("NewPrescription", nil))
}

// @Summary Delete a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionID path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID} [DELETE]
func (h *CaseHandler) DeletePrescription(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("DeletePrescription", nil))
}

// @Summary Update a prescrition
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionDetails body uint true "prescription ID and updated details"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID} [PUT]
func (h *CaseHandler) UpdatePrescription(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("UpdatePrescription", nil))
}

// @Summary Get prescription by prescription id
// @Description
// @Tags Case
// @Produce json
// @Param prescriptionID path uint true "prescription ID"
// @Success 200 {object} api.ReturnedData{data=Prescription}
// @Router /patient/{patientID}/case/{caseID}/prescription/{prescriptionID}  [GET]
func (h *CaseHandler) GetPrescriptionByPrescriptionID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPrescriptionByPrescriptionID", nil))
}

// @Summary Get prescriptions by case id
// @Description
// @Tags Case
// @Produce json
// @Param caseID path uint true "case ID"
// @Success 200 {object} api.ReturnedData{data=[]Prescription}
// @Router /patient/{patientID}/case/{caseID}/prescription  [GET]
func (h *CaseHandler) GetPrescriptionByCaseID(c echo.Context) error {
	// ...
	c.Logger().Debug("hello world")
	return c.JSON(200, api.Return("GetPrescriptionByCaseID", nil))
}
