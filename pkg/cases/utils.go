package cases

import (
	"github.com/AsterNighT/software-engineering-backend/pkg/database/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func FromAdmin(c echo.Context) bool {
	db, _ := c.Get("db").(*gorm.DB)
	id := c.Get("id").(uint)
	var user models.Account
	if err := db.Where("id = ?", id).First(&user).Error; err == nil { // not found
		return user.Type == models.AdminType
	}
	return false
}

func FromDoctor(c echo.Context) bool {
	if FromAdmin(c) {
		c.Logger().Debug("Accept Admin")
		return true
	}
	db, _ := c.Get("db").(*gorm.DB)
	id := c.Get("id").(uint)
	var user models.Account
	if err := db.Where("id = ?", id).First(&user).Error; err == nil { // not found
		return user.Type == models.DoctorType
	}
	return false
}

func FromPatient(c echo.Context, id uint) bool {
	if FromDoctor(c) {
		c.Logger().Debug("Accept Doctor")
		return true
	}
	if c.Get("id").(uint) == id {
		c.Logger().Debug("Accept Patient", id)
		return true
	}
	c.Logger().Errorf("unauthorized access, expecting user with id:%s", id)
	return false
}
