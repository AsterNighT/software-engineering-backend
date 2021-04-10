package account

import (
	"container/list"
	"crypto/rand"
	"net/http"
	"regexp"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/labstack/echo/v4"
)

var accountList = list.New()

type AccountHandler struct {
}

// @Summary create and account based on email(as id), type, name and password
// @Description will check primarykey other, then add to accountList if possible
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Type path string true "user type"
// @Param Name path string true "user name"
// @Param Passwd path string true "user password"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/account_table [POST]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	Email := c.QueryParam("Email")
	Type := AcountType(c.QueryParam("Type"))
	Name := c.QueryParam("Name")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invali E-mail Address", nil))
	}

	if Type != patient && Type != doctor && Type != admin {
		return c.JSON(http.StatusBadRequest, api.Return("Invali Account Type", nil))
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, api.Return("Invali Password Length", nil))
	}

	// Check uniqueness
	for itor := accountList.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			return c.JSON(http.StatusBadRequest, api.Return("E-Mail occupied", nil))
		}
	}

	accountList.PushBack(Account{Email: Email, Type: Type, Name: Name, Passwd: Passwd})

	return c.JSON(http.StatusOK, api.Return("Successfully created", nil))

}

/**
 * @todo cookie not implemented, jwt
 */

// @Summary login using email and passwd
// @Description
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Passwd path string true "user password"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account [POST]
func (h *AccountHandler) LoginAccount(c echo.Context) error {
	Email := c.QueryParam("Email")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invali E-mail Address", nil))
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, api.Return("Invali Password Length", nil))
	}

	return checkPasswd(c, Email, Passwd)
}

// @Summary reset password
// @Description host will send a verification code to email, need response with verification code
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @param VeriCode path string true "verification code sent by user"
// @Param Passwd path string true "user password"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/{id}/reset [PUT]
func (h *AccountHandler) ResetPasswd(c echo.Context) error {
	Email := c.QueryParam("Email")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invali E-mail Address", nil))
	}

	// Gen verification code
	buffer := make([]byte, 6)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	for i := 0; i < 6; i++ {
		buffer[i] = "1234567890"[int(buffer[i])%6]
	}
	// hostVcode := string(buffer)

	// SendVeriMsg(Email, hostVcode) // Func wait for implementation

	// Wait for response from client...

	clientVcode := c.QueryParam("VeriCode")
	newPasswd := c.QueryParam("Passwd")

	// if clientVcode == hostVcode {
	if clientVcode == string(buffer) {
		return modifyPasswd(c, Email, newPasswd)
	} else {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Verification Code", nil))
	}
}

// @Summary the interface of modifying password
// @Description can only be called during logged-in status since there is no password check
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Passwd path string true "user password (the new one)"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/{id}/modify [PUT]
func (h *AccountHandler) ModifyPasswd(c echo.Context) error {
	Email := c.QueryParam("Email")
	Passwd := c.QueryParam("Passwd")

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invali E-mail Address", nil))
	}

	if len(Passwd) < ACCOUNT_PASSWD_LEN {
		return c.JSON(http.StatusBadRequest, api.Return("Invali Password Length", nil))
	}
	return modifyPasswd(c, Email, Passwd)
}

/**
 * @brief private method for checking password
 */
func checkPasswd(c echo.Context, Email string, Passwd string) error {
	// Travese to find matched account, use DB later
	for itor := accountList.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			if itor.Value.(Account).Passwd == Passwd {
				return c.JSON(http.StatusOK, api.Return("Successfully logged in", nil))
			} else {
				return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
			}
		}
	}

	return c.JSON(http.StatusBadRequest, api.Return("E-Mail not found", nil))
}

/**
 * @brief private method for modifying password
 */
func modifyPasswd(c echo.Context, Email string, Passwd string) error {
	for itor := accountList.Front(); itor != nil; itor = itor.Next() {
		if itor.Value.(Account).Email == Email {
			// Remove this and push a new one with new passwd
			accountList.PushBack(Account{Email: Email, Type: itor.Value.(Account).Type, Name: itor.Value.(Account).Name, Passwd: Passwd})
			accountList.Remove(itor)

			return c.JSON(http.StatusOK, api.Return("Successfully modified", nil))
		}
	}
	return c.JSON(http.StatusBadRequest, api.Return("E-Mail not found", nil))
}
