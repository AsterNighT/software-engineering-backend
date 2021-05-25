package utils

import "gorm.io/gorm"

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}
