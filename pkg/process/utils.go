package process

import (
	"time"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func InitProcessValidator() error {
	validate = validator.New()

	err := validate.RegisterValidation("halfday", validateHalfDay)
	if err != nil {
		return err
	}
	return nil
}

func validateHalfDay(fl validator.FieldLevel) bool {
	s := HalfDayEnum(fl.Field().String())
	return s == morning || s == afternoon
}

func validateSchedule(schedule *DepartmentSchedule) bool {
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
