package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
	"github.com/AsterNighT/software-engineering-backend/pkg/cases"
	"github.com/AsterNighT/software-engineering-backend/pkg/process"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func InitDb() *gorm.DB {
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

	// auto migrate cases
	err = db.AutoMigrate(&cases.Prescription{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&cases.Case{})
	if err != nil {
		panic(err)
	}

	// auto migrate account
	err = db.AutoMigrate(
		&account.Account{},
		&account.Patient{},
		&account.Doctor{},
	)

	if err != nil {
		panic(err)
	}

	return db
}

func ContextDB(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}
