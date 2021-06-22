package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
	"github.com/AsterNighT/software-engineering-backend/pkg/cases"
	"github.com/AsterNighT/software-engineering-backend/pkg/process"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func InitDB() *gorm.DB {
	dsn := "host=neon-cubes.xyz user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&cases.Medicine{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&cases.Guideline{})
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
	err = process.InitProcessValidator()
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

	utils.DB = db

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
