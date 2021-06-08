package process

import (
	"time"

	"github.com/AsterNighT/software-engineering-backend/pkg/account"
)

// Department
//department table (e.g. orthopedics department, x-ray department)
type Department struct {
	ID        uint                 `json:"id" gorm:"primaryKey"`
	Name      string               `json:"name"`                   // name of this department
	Detail    string               `json:"detail"`                 // detailed introduction of this department
	Doctors   []account.Doctor     `json:"-" swaggerignore:"true"` // foreign key of all the doctors who belongs to this department
	Schedules []DepartmentSchedule `json:"-" swaggerignore:"true"` // time schedule for a whole department
}

// Registration
// registration table
// every registration will eventually be terminated, and therefore needs a cause
type Registration struct {
	ID              uint            `gorm:"primaryKey"`
	DoctorID        uint            `swaggerignore:"true"`
	Doctor          account.Doctor  `swaggerignore:"true"`
	PatientID       uint            `swaggerignore:"true"`
	Patient         account.Patient `swaggerignore:"true"`
	DepartmentID    uint            `swaggerignore:"true"`
	Department      Department      `swaggerignore:"true"`
	Date            time.Time
	HalfDay         HalfDayEnum            // TODO: a validator for registration, only half day is allowed
	Status          RegistrationStatusEnum `gorm:"default:'committed'"`
	TerminatedCause string                 `gorm:"default''"`
	MileStones      []MileStone
}

// MileStone
// milestone that represent a small step during the process
type MileStone struct {
	ID             uint         `gorm:"primaryKey" swaggerignore:"true"`
	RegistrationID uint         `swaggerignore:"true"`
	Registration   Registration `swaggerignore:"true"`
	Activity       string       `gorm:"default:''"`
	Checked        bool         `gorm:"default:false"`
}

// DepartmentSchedule
// schedule table for a whole department
// each object represents a minimal schedule duration
type DepartmentSchedule struct {
	ID           uint       `gorm:"primaryKey" swaggerignore:"true"`
	DepartmentID uint       `swaggerignore:"true"`
	Department   Department `swaggerignore:"true"`
	Date         time.Time
	HalfDay      HalfDayEnum // TODO: a validator for department, only half day is allowed
	Capacity     int
	Current      int // current number of registrations of this schedule duration
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
	whole     HalfDayEnum = "whole"
)
