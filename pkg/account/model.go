package account

// "github.com/AsterNighT/software-engineering-backend/pkg/cases"
// "github.com/AsterNighT/software-engineering-backend/pkg/chat"

const accountPasswdLen = 8

type Account struct {
	ID    uint `gorm:"primarykey;autoIncrement;"`
	Email string

	Type AcountType
	// Name   string
	FirstName string
	LastName  string
	Passwd    string // Wait for encryption

	Token string
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
	jwtKey = []byte("ECHOGORMEXAMPLEJWTKEY") // For test only
	// jwtKey = os.Getenv("JWT_KEY")
)
