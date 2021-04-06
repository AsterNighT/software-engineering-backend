package process

import (
	"time"
)

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// these tables should be included in account model (maybe)
type Doctor struct {
	/* ... */
	Department Department
}

type Patient struct {
	/* ... */
}

type DoctorSchedule struct {
	ID       uint `gorm:"primaryKey"`
	Doctor   Doctor
	Date     time.Time
	HalfDay  HalfDayEnum `gorm:"default:'whole'"`
	Capacity int
}

//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// department table (e.g. orthopedics department, x-ray department)
type Department struct {
	ID        uint                 `gorm:"primaryKey"`
	Name      string               // name of this department
	Detail    string               // detailed introduction of this department
	Doctors   []Doctor             // foreign key of all the doctors who belongs to this department
	Schedules []DepartmentSchedule // time schedule for a whole department
}

// registration table
type Registration struct {
	ID         uint `gorm:"primaryKey"`
	Doctor     Doctor
	Patient    Patient
	Department Department
	Date       time.Time
	HalfDay    HalfDayEnum            // TODO: a validator for registration, only half day is allowed
	Status     RegistrationStatusEnum `gorm:"default:'committed'"`
	// every registration will eventually be terminated, and therefore needs a cause
	TerminatedCause string `gorm:"default''"`
	MileStones      []MileStone
}

// milestone that represent a small step during the process
type MileStone struct {
	ID           uint `gorm:"primaryKey"`
	Time         time.Time
	Registration Registration

	Location string `gorm:"default:''"`
	Activity string `gorm:"default:''"`
	Assignee string `gorm:"default:''"` // someone who is responsible for this small milestone
	Checked  bool   `gorm:"default:false"`
}

// schedule table for a whole department
// each object represents a minimal schedule duration
type DepartmentSchedule struct {
	ID         uint `gorm:"primaryKey"`
	Department Department
	Date       time.Time
	HalfDay    HalfDayEnum // TODO: a validator for department, only half day is allowed
	Capacity   int
	// DepartmentSchedule.Capacity = SUM(DoctorSchedule.Capacity if the doctor belongs to this department)
	Current   int // current number of registrations of this schedule duration
}

// define new enum for registration status
type RegistrationStatusEnum string

const (
	committed  RegistrationStatusEnum = "committed"
	accepted   RegistrationStatusEnum = "accepted"
	terminated RegistrationStatusEnum = "terminated"
)

// define new enum for half day selection
type HalfDayEnum string

const (
	morning   HalfDayEnum = "morning"
	afternoon HalfDayEnum = "afternoon"
	whole     HalfDayEnum = "whole"
)
