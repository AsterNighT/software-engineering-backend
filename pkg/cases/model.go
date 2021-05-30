package cases

import (
	"time"
)

type Case struct {
	ID             uint `gorm:"primaryKey"`   // Every object should have ID
	PatientID      uint `validate:"required"` // A has many relationship should be on this
	DoctorID       uint
	Department     string    `validate:"required"`
	Complaint      string    `validate:"required"` // Use urls to locate pictures
	Diagnosis      string    `validate:"required"`
	Treatment      string    `validate:"required"`
	History        string    `validate:"required"`
	Date           time.Time `validate:"required,datetime"`
	Prescriptions  []Prescription
	PreviousCase   *Case // Previous case (the lastest one). If there is none previous case, set nil
	PreviousCaseID *uint `validate:"required"`
}

type Prescription struct {
	ID         uint `gorm:"primarykey" json:"omitempty"`
	CaseID     uint
	Advice     string `validate:"required"`
	Guidelines []Guideline
}

type Guideline struct {
	ID             uint `gorm:"primarykey"`
	MedicineID     uint `validate:"required"`
	Medicine       Medicine
	PrescriptionID uint
	Dosage         string `validate:"required"`
	Quantity       uint   `validate:"required"`
}

type Medicine struct {
	ID               uint `gorm:"primarykey"`
	Name             string `validate:"required"`
	Price            float32 `validate:"required"`
	Contraindication string `validate:"required"`
}
