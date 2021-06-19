package account

import "time"

const accountPasswdLen = 8

type Account struct {
	ID    uint `gorm:"primarykey;autoIncrement;"`
	Email string

	Type      AcountType
	FirstName string
	LastName  string
	Passwd    string // Considered as plaintext, but can be encrypted by frontend

	Token           string
	AuthCode        string
	AuthCodeExpires time.Time
}

type AcountType string

const (
	PatientType AcountType = "patient"
	DoctorType  AcountType = "doctor"
	AdminType   AcountType = "admin"
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
	jwtKey = []byte("SOFTWAREENGINEERINGBACKEND") // For test only

	emailServerHost = "smtp.163.com"
	emailServerPort = "25"
	emailUser       = "13606900243@163.com"
	emailPasswd     = "IAXPDCMBTUCJMXOC"

	expireMin = 10
)
