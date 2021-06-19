package process

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
)

// Department
// department table (e.g. orthopedics department, x-ray department)
type Department struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`   // name of this department
	Detail string `json:"detail"` // detailed introduction of this department
}

// Registration
// registration table
// every registration will eventually be terminated, and therefore needs a cause
type Registration struct {
	ID           uint            `gorm:"primaryKey"`
	DoctorID     uint            `swaggerignore:"true"`
	Doctor       account.Doctor  `swaggerignore:"true"`
	PatientID    uint            `swaggerignore:"true"`
	Patient      account.Patient `swaggerignore:"true"`
	DepartmentID uint            `swaggerignore:"true"`
	Department   Department      `swaggerignore:"true"`

	Year    int         `json:"year"`
	Month   int         `json:"month"`
	Day     int         `json:"day"`
	HalfDay HalfDayEnum `json:"halfday"`

	Status          RegistrationStatusEnum `gorm:"default:'committed'"`
	TerminatedCause string                 `gorm:"default''"`
}

// MileStone
// milestone that represent a small step during the process
type MileStone struct {
	ID             uint         `json:"id,omitempty" gorm:"primaryKey" swaggerignore:"true"`
	RegistrationID uint         `json:"-" swaggerignore:"true"`
	Registration   Registration `json:"-" swaggerignore:"true"`
	Activity       string       `json:"activity,omitempty" gorm:"default:''"`
	Checked        bool         `json:"checked,omitempty" gorm:"default:false"`
}

// DepartmentSchedule
// schedule table for a whole department
// each object represents a minimal schedule duration
type DepartmentSchedule struct {
	ID           uint        `json:"-" gorm:"primaryKey" swaggerignore:"true"`
	DepartmentID uint        `json:"-" swaggerignore:"true"`
	Department   Department  `json:"-" swaggerignore:"true"`
	Year         int         `json:"year"`
	Month        int         `json:"month"`
	Day          int         `json:"day"`
	HalfDay      HalfDayEnum `json:"halfday"`
	Capacity     int         `json:"capacity"`
	Current      int         `json:"current"` // current number of registrations of this schedule duration
}

// RegistrationStatusEnum
// define new enum for registration status
type RegistrationStatusEnum string

const (
	committed  RegistrationStatusEnum = "committed"
	accepted   RegistrationStatusEnum = "accepted"
	terminated RegistrationStatusEnum = "terminated"
)

// HalfDayEnum
// define new enum for half day selection
type HalfDayEnum string

const (
	morning   HalfDayEnum = "morning"
	afternoon HalfDayEnum = "afternoon"
)

type ProcessError string

const (
	InvalidSubmitFormat   ProcessError = "invalid submit format"
	DepartmentNotFound    ProcessError = "department not found"
	PatientNotFound       ProcessError = "patient not found"
	DoctorNotFound        ProcessError = "doctor not found"
	DuplicateRegistration ProcessError = "duplicate registration is not allowed"
	InvalidSchedule       ProcessError = "this schedule is invalid"
	NotEnoughCapacity     ProcessError = "not enough capacity"
	CannotAssignDoctor    ProcessError = ""
	InvalidRegistration   ProcessError = ""
)

// DepartmentDetailJSON defines the return json to frontend
type DepartmentDetailJSON struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Detail    string               `json:"detail"`
	Doctors   []string             `json:"doctors"`
	Schedules []DepartmentSchedule `json:"schedules"`
}

type RegistrationJSON struct {
	ID         uint                   `json:"id"`
	Department string                 `json:"department"`
	Status     RegistrationStatusEnum `json:"status"`
	Year       int                    `json:"year"`
	Month      int                    `json:"month"`
	Day        int                    `json:"day"`
	HalfDay    HalfDayEnum            `json:"halfday"`
}

type RegistrationDetailJSON struct {
	ID         uint                   `json:"id"`
	Department string                 `json:"department"`
	Doctor     string                 `json:"doctor"`
	Patient    string                 `json:"patient"`
	Year       int                    `json:"year"`
	Month      int                    `json:"month"`
	Day        int                    `json:"day"`
	HalfDay    HalfDayEnum            `json:"halfday"`
	Status     RegistrationStatusEnum `json:"status"`
	MileStone  MileStone              `json:"milestone,omitempty"`
}
