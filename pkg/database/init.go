package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/cases"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
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
	err = db.AutoMigrate(&cases.Prescription{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&cases.Case{})
	if err != nil {
		panic(err)
	}
	utils.DB = db
}
