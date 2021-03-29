package account

const ACCOUNT_PASSWD_LEN = 8

type Account struct {
	ID    uint `gorm:"primarykey"`
	Email string

	Type   AcountType
	Name   string
	Passwd string // Wait for encryption
}

type AcountType string

const (
	ACCOUNT_TYPE_PATIENT AcountType = "patient"
	ACCOUNT_TYPE_DOCTOR  AcountType = "doctor"
	ACCOUNT_TYPE_ADMIN   AcountType = "admin"
)
