package cases

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/account"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func FromAdmin(c echo.Context) bool {
	db, _ := c.Get("db").(*gorm.DB)
	id := c.Get("id").(string)
	var user account.Account
	if err := db.Where("id = ?", id).First(&user).Error; err == nil { // not found
		return user.Type == account.AdminType
	}
	return false
}

func FromDoctor(c echo.Context) bool {
	if FromAdmin(c) {
		return true
	}
	db, _ := c.Get("db").(*gorm.DB)
	id := c.Get("id").(string)
	var user account.Account
	if err := db.Where("id = ?", id).First(&user).Error; err == nil { // not found
		return user.Type == account.DoctorType
	}
	return false
}

func FromPatient(c echo.Context, id string) bool {
	if FromDoctor(c) {
		return true
	}
	if c.Get("id").(string) == id {
		return true
	}
	c.Logger().Errorf("unauthorized access, expecting user with id:%s", id)
	return false
}
