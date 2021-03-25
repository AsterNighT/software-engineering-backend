package cases

type Case struct {
	ID            uint `gorm:"primarykey"` // Every object should have ID
	PatientID     uint // A has many relationship should be on this
	Complaint     string
	Diagnosis     string
	Prescriptions []Prescription
}

type Prescription struct {
	ID      uint `gorm:"primarykey"`
	CaseID  uint
	Details string
}
