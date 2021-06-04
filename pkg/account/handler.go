package account

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// var accountList = list.New()

type AccountHandler struct {
}

// Init : Init Router
// func (h AccountHandler) Init(g *echo.Group) {
// 	g.POST("/create", h.CreateAccount)
// 	g.POST("/login", h.LoginAccount)
// 	g.POST("/logout", h.LogoutAccount, Authoriszed)
// }

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
	type RequestBody struct {
		ID    uint   `json:"id" validate:"required"`
		Email string `json:"email" validate:"required"`

		Type   AcountType `json:"type" validate:"required"`
		Name   string     `json:"name" validate:"required"`
		Passwd string     `json:"passwd" validate:"required"`
	}

	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	if body.Type != patient && body.Type != doctor && body.Type != admin {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Account Type", nil))
	}

	if len(body.Passwd) < accountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	if err := db.Where("id = ? OR email = ?", body.ID, body.Email).First(&Account{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail or AccountID occupied", nil))
	}

	account := Account(body)

	// accountList.PushBack(Account{Email: Email, Type: Type, Name: Name, Passwd: Passwd})
	account.HashPassword()

	if result := db.Create(&account); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error))
	}

	token, _ := account.GenerateToken()
	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	c.SetCookie(&cookie)

	// return c.JSON(http.StatusOK, api.Return("Successfully created", nil))
	return c.JSON(http.StatusOK, api.Return("Created", echo.Map{
		"account":      account,
		"cookie_token": token,
	}))
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
	type RequestBody struct {
		Email  string `json:"email" validate:"required"`
		Passwd string `json:"passwd" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}
	if len(body.Passwd) < accountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	var account Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail not found", nil))
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(body.Passwd)) != nil {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
	}

	token, _ := account.GenerateToken()
	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, api.Return("Logged in", echo.Map{
		"account":      account,
		"cookie_token": token,
	}))
}

// @Summary logout using cookie
// @Description
// @Tags Account
// @Produce json
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account [GET]
func (h *AccountHandler) LogoutAccount(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusBadRequest, api.Return("Not Logged in", nil))
	}
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, api.Return("Account logged out", nil))
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
// func (h *AccountHandler) ResetPasswd(c echo.Context) error {
// 	Email := c.QueryParam("Email")

// 	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, Email); !ok {
// 		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
// 	}

// 	// Gen verification code
// 	buffer := make([]byte, 6)
// 	if _, err := rand.Read(buffer); err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < 6; i++ {
// 		buffer[i] = "1234567890"[int(buffer[i])%6]
// 	}
// 	// hostVcode := string(buffer)

// 	// SendVeriMsg(Email, hostVcode) // Func wait for implementation

// 	// Wait for response from client...

// 	clientVcode := c.QueryParam("VeriCode")
// 	newPasswd := c.QueryParam("Passwd")

// 	// if clientVcode == hostVcode {
// 	if clientVcode == string(buffer) {
// 		return modifyPasswd(c, Email, newPasswd)
// 	}
// 	return c.JSON(http.StatusBadRequest, api.Return("Wrong Verification Code", nil))

// }

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
	type RequestBody struct {
		Email     string `json:"email" validate:"required"`
		Passwd    string `json:"passwd" validate:"required"`
		NewPasswd string `json:"newpasswd" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	// Check old passwd
	db, _ := c.Get("db").(*gorm.DB)
	var account Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail not found", nil))
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(body.Passwd)) != nil {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
	}

	if len(body.NewPasswd) < accountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	account.Passwd = body.NewPasswd
	account.HashPassword()

	if result := db.Model(&Account{}).Where("id = ?", account.ID).Update("passwd", account.Passwd); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error))
	}

	return c.JSON(http.StatusOK, api.Return("Successfully modified", nil))
}

/**
 * @brief private method for hashing password
 */
func (u *Account) HashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(u.Passwd), bcrypt.DefaultCost)
	u.Passwd = string(bytes)
}

/**
 * @brief private method for generateing cookie token
 */
func (u *Account) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.ID,
	})
	// print(token)
	tokenString, err := token.SignedString(jwtKey)
	// print(err.Error())
	return tokenString, err
}

// Authoriszed : Check Auth
func Authoriszed(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})

		if !token.Valid || err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("id", token.Claims.(jwt.MapClaims)["id"])

		return next(c)
	}
}
