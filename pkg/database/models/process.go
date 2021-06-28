package models

import (
	"time"

	"github.com/go-playground/validator"
)

// Department
// department table (e.g. orthopedics department, x-ray department)
type Department struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"unique"` // name of this department
	Detail    string `json:"detail"`             // detailed introduction of this department
	Questions string `json:"questions"`
}

// Registration
// registration table
// every registration will eventually be terminated, and therefore needs a cause
type Registration struct {
	ID           uint       `gorm:"primaryKey"`
	DoctorID     uint       `swaggerignore:"true"`
	Doctor       Doctor     `swaggerignore:"true"`
	PatientID    uint       `swaggerignore:"true"`
	Patient      Patient    `swaggerignore:"true"`
	DepartmentID uint       `swaggerignore:"true"`
	Department   Department `swaggerignore:"true"`

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
	Checked        bool         `json:"checked" gorm:"default:false"`
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
	Committed  RegistrationStatusEnum = "committed"
	Accepted   RegistrationStatusEnum = "accepted"
	Terminated RegistrationStatusEnum = "terminated"
)

// HalfDayEnum
// define new enum for half day selection
type HalfDayEnum string

const (
	Morning   HalfDayEnum = "morning"
	Afternoon HalfDayEnum = "afternoon"
)

type ProcessError string

const (
	InvalidSubmitFormat      ProcessError = "参数格式错误"
	DepartmentNotFound       ProcessError = "无法找到该科室"
	PatientNotFound          ProcessError = "无法找到该患者"
	RegistrationNotFound     ProcessError = "找不到该挂号"
	DoctorNotFound           ProcessError = "无法找到该医生"
	MileStoneNotFound        ProcessError = "找不到该 MileStone"
	MileStoneUnauthorized    ProcessError = "你无权操作该 MileStone"
	InvalidSchedule          ProcessError = "该时段不可用"
	NotEnoughCapacity        ProcessError = "该时段已经没有足够的挂号余量"
	CreateRegistrationFailed ProcessError = "挂号失败，请检查科室是否有余量或者时间冲突"
	CreateMileStoneFailed    ProcessError = "创建 MileStone 失败"
	CannotAssignDoctor       ProcessError = "无法为该挂号分配医生"
	AccountNotFound          ProcessError = "找不到用户"
	RegistrationUpdateFailed ProcessError = "挂号状态更新失败"
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
	ID              uint                   `json:"id"`
	Department      string                 `json:"department"`
	Doctor          string                 `json:"doctor"`
	Patient         string                 `json:"patient"`
	Year            int                    `json:"year"`
	Month           int                    `json:"month"`
	Day             int                    `json:"day"`
	HalfDay         HalfDayEnum            `json:"halfday"`
	Status          RegistrationStatusEnum `json:"status"`
	MileStone       []MileStone            `json:"milestone,omitempty"`
	TerminatedCause string                 `json:"terminated_cause"`
}

var Validate *validator.Validate

func InitProcessValidator() error {
	Validate = validator.New()

	err := Validate.RegisterValidation("halfday", ValidateHalfDay)
	if err != nil {
		return err
	}
	return nil
}

func ValidateHalfDay(fl validator.FieldLevel) bool {
	s := HalfDayEnum(fl.Field().String())
	return s == Morning || s == Afternoon
}

func ValidateSchedule(schedule *DepartmentSchedule) bool {
	thisYear, thisMonth, thisDay := time.Now().Date()
	if schedule.Year < thisYear {
		return false
	} else if schedule.Year == thisYear && schedule.Month < int(thisMonth) {
		return false
	} else if schedule.Year == thisYear && schedule.Month == int(thisMonth) && schedule.Day <= thisDay {
		return false
	}

	return true
}
