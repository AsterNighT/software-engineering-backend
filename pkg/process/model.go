package process

import (
	"gorm.io/gorm"
	"time"
)

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// these tables should be included in account model (maybe)
type Doctor struct {
	gorm.Model
	/* ... */
	Department Department
}

type Patient struct {
	gorm.Model
	/* ... */
}

type DoctorSchedule struct {
	gorm.Model
	Doctor   Doctor
	Date     time.Time
	HalfDay  HalfDayEnum `gorm:"default:0"`
	Capacity int
}

//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// department table (e.g. orthopedics department, x-ray department)
type Department struct {
	gorm.Model
	Name      string               // name of this department
	Detail    string               // detailed introduction of this department
	Doctors   []Doctor             // foreign key of all the doctors who belongs to this department
	Schedules []DepartmentSchedule // time schedule for a whole department
}

// registration table
type Registration struct {
	gorm.Model
	Doctor     Doctor
	Patient    Patient
	Department Department

	Status RegistrationStatusEnum `gorm:"default:2"`
	// every registration will eventually be terminated, and therefore needs a cause
	TerminatedCause string `gorm:"default''"`
	MileStones      []MileStone
}

// milestone that represent a small step during the process
type MileStone struct {
	gorm.Model
	Time         time.Time
	Registration Registration

	Location string `gorm:"default:''"`
	Activity string `gorm:"default:''"`
	Assignee string `gorm:"default:''"` // someone who is responsible for this small milestone
	Checked  bool   `gorm:"default:false"`
}

// schedule table for a whole department
type DepartmentSchedule struct {
	gorm.Model
	Department Department
	Date       time.Time
	HalfDay    HalfDayEnum `gorm:"default:0"`
	Capacity   int
	// DepartmentSchedule.Capacity = SUM(DoctorSchedule.Capacity if the doctor belongs to this department)
}

// define new enum for registration status
type RegistrationStatusEnum int

const (
	committed  RegistrationStatusEnum = 0
	accepted   RegistrationStatusEnum = 1
	terminated RegistrationStatusEnum = 2
)

// define new enum for half day selection
type HalfDayEnum int

const (
	morning   HalfDayEnum = 0
	afternoon HalfDayEnum = 1
	full      HalfDayEnum = 2
)
