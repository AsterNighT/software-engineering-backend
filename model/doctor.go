package model

import (
	"github.com/go-gorm/gorm"
)

type Doctor struct {
	gorm.Model
	DepartmentId uint
	//TODO: 剩余属性其他相关的组完善

}
