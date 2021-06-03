package account

import (
	// "github.com/AsterNighT/software-engineering-backend/pkg/cases"
	// "github.com/AsterNighT/software-engineering-backend/pkg/chat"
)

const accountPasswdLen = 8

type Account struct {
	ID    uint `gorm:"primarykey;"`
	Email string

	Type   AcountType
	Name   string
	Passwd string // Wait for encryption
}

type AcountType string

const (
	patient AcountType = "patient"
	doctor  AcountType = "doctor"
	admin   AcountType = "admin"
)

type Doctor struct {
	ID           uint `gorm:"primarykey"`
	DepartmentID uint

	AccountID uint
	CaseID    uint
	ChatID    uint
	// Cases     []cases.Case `gorm:"foreignkey:ID"`
	// Chats     []chat.Chat  `gorm:"foreignkey:ID"`
}

type Patient struct {
	ID uint `gorm:"primarykey"`

	AccountID uint
	CaseID    uint
	ChatID    uint
	// Account Account      `gorm:"foreignkey:ID"`
	// Cases   []cases.Case `gorm:"foreignkey:ID"`
	// Chats   []chat.Chat  `gorm:"foreignkey:ID"`
}

var (
	jwtKey = []byte("ECHOGORMEXAMPLEJWTKEY") // For test only
	// jwtKey = os.Getenv("JWT_KEY")
)
