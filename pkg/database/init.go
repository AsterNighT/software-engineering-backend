package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/database/models"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Medicine{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Guideline{})
	if err != nil {
		panic(err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//		G4 - Process
	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	// auto migrate process's table
	err = db.AutoMigrate(
		&models.Department{},
		&models.Registration{},
		&models.MileStone{},
		&models.DepartmentSchedule{},
	)
	if err != nil {
		panic(err)
	}
	err = models.InitProcessValidator()
	if err != nil {
		panic(err)
	}

	// auto migrate cases
	err = db.AutoMigrate(&models.Prescription{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Case{})
	if err != nil {
		panic(err)
	}

	utils.DB = db

	// auto migrate account
	err = db.AutoMigrate(
		&models.Account{},
		&models.Auth{},
		&models.Patient{},
		&models.Doctor{},
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
