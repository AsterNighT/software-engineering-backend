package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/process"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//		G4 - Process
	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	// auto migrate process's table
	err = db.AutoMigrate(
		&process.Department{},
		&process.Registration{},
		&process.MileStone{},
		&process.DepartmentSchedule{},
	)

	if err != nil {
		panic(err)
	}

	//err = db.AutoMigrate(&cases.Prescription{})
	//if err != nil {
	//	panic(err)
	//}
	//err = db.AutoMigrate(&cases.Case{})
	//if err != nil {
	//	panic(err)
	//}
}
