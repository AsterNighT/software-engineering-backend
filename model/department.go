package model

import (
	"github.com/go-gorm/gorm"
)

type Department struct {
	gorm.Model
	DepartmentName string 		//科室名
	Detail string				//科室简介
	Capacity int				//科室容量，
	Remain int
	Waiting int //Waiting = Capacity - Remain
}
