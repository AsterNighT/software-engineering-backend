package account

const ACCOUNT_PASSWD_LEN = 8

type Account struct {
	Email string `gorm:"primarykey"`

	// Type   AcountType
	Type   int
	Name   string
	Passwd string // Wait for encryption
}

type AcountType int

const (
	ACCOUNT_TYPE_PATIENT AcountType = iota
	ACCOUNT_TYPE_DOCTOR  AcountType = iota
	ACCOUNT_TYPE_ADMIN   AcountType = iota
)
