package database

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/database/models"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

	initDepartment()
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

// initDepartment initialize the department tables.
func initDepartment() {
	departments := []models.Department{{
		Name:   "眼科",
		Detail: "世界一流眼科",
	}, {
		Name:   "骨科",
		Detail: "",
	}, {
		Name:   "呼吸内科",
		Detail: "",
	}, {
		Name:   "妇产科",
		Detail: "专治不孕不育",
	}, {
		Name:   "肿瘤内科",
		Detail: "本科主要针对消化系统癌症及肺癌等血液恶性肿瘤以外的固态瘤病人进行诊疗（包括诊断、癌症药物治疗、与放射科及外科共同进行多学科综合治疗",
	}, {
		Name:   "心胸外科",
		Detail: "",
	}, {
		Name:   "耳鼻喉科",
		Detail: "",
	}, {
		Name:   "急诊科",
		Detail: "",
	}}

	db := utils.GetDB()
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"detail"}),
	}).Create(&departments)
}
