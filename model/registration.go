package model

import (
	"github.com/go-gorm/gorm"
	"time"
)

type Registration struct {
	gorm.Model
	PatientId uint					//
	DepartmentId uint				//科室id，传递依赖
	DoctorId uint
	RegistrationTime time.Time		//挂号时间，挂号完成之后（进行余量判断），创建空白订单

}

