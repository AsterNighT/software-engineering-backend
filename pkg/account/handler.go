package account

import (
	"container/list"
	"crypto/rand"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

var account_list = list.New()

// type AcountType string

type AccountHandler struct {
}

/**
 * @brief create and account based on email(as id), type, name and password
 * @desc will check primarykey other, then add to account_list if possible
 */
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	Email := c.QueryParam("Email")
	Type := AcountType(c.QueryParam("Type"))
	Name := c.QueryParam("Name")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali E-mail Address"})
	}

	if Type != ACCOUNT_TYPE_PATIENT && Type != ACCOUNT_TYPE_DOCTOR && Type != ACCOUNT_TYPE_ADMIN {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali Account Type"})
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali Password Length"})
	}

	// Check uniqueness
	for itor := account_list.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "E-Mail occupied"})
		}
	}

	account_list.PushBack(Account{Email: Email, Type: Type, Name: Name, Passwd: Passwd})

	return c.String(http.StatusOK, "Successfully created")

}

/**
 * @brief login using email and passwd
 * @desc
 * @todo cookie not implemented
 */
func (h *AccountHandler) LoginAccount(c echo.Context) error {
	Email := c.QueryParam("Email")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali E-mail Address"})
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali Password Length"})
	}

	return checkPasswd(c, Email, Passwd)
}

/**
 * @brief reset password
 * @desc host will send a verification code to email, need response with verification code
 */
func (h *AccountHandler) ResetPasswd(c echo.Context) error {
	Email := c.QueryParam("Email")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali E-mail Address"})
	}

	// Gen verification code
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 6; i++ {
		buffer[i] = "1234567890"[int(buffer[i])%6]
	}
	host_vcode := string(buffer)

	// SendVeriMsg(Email, host_vcode) // Func wait for implementation

	// Wait for response from client...

	client_vcode := c.QueryParam("VeriCode")
	new_passwd := c.QueryParam("Passwd")

	if client_vcode == host_vcode {
		return modifyPasswd(c, Email, new_passwd)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Wrong Verification Code"})
	}
}

/**
 * @brief the interface of modifying password
 * @note can only be called during logged-in status since there is no password check
 */
func (h *AccountHandler) ModifyPasswd(c echo.Context) error {
	Email := c.QueryParam("Email")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali E-mail Address"})
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invali Password Length"})
	}
	return modifyPasswd(c, Email, Passwd)
}

/**
 * @brief private method for checking password
 */
func checkPasswd(c echo.Context, Email string, Passwd string) error {
	// Travese to find matched account, use DB later
	for itor := account_list.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			if itor.Value.(Account).Passwd == Passwd {
				return c.String(http.StatusOK, "Successfully Logged in")
			} else {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Wrong Password"})
			}
		}
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "E-Mail not found"})
}

/**
 * @brief private method for modifying password
 */
func modifyPasswd(c echo.Context, Email string, Passwd string) error {
	for itor := account_list.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			// Remove this and push a new one with new passwd
			account_list.PushBack(Account{Email: Email, Type: itor.Value.(Account).Type, Name: itor.Value.(Account).Name, Passwd: Passwd})
			account_list.Remove(itor)

			return c.String(http.StatusOK, "Successfully modified")
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "E-Mail not found"})
}
