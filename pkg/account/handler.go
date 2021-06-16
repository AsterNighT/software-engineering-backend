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

type AccountHandler struct{}

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
// @Router /account/create [POST]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	type RequestBody struct {
		ID    string `json:"id" validate:"required"`
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
	if body.Type != PatientType && body.Type != DoctorType && body.Type != AdminType {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Account Type", nil))
	}
	if len(body.Passwd) < accountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	if err := db.Where("id = ? OR email = ?", body.ID, body.Email).First(&Account{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail or AccountID occupied", nil))
	}

	account := Account{
		ID:    body.ID,
		Email: body.Email,

		Type:   body.Type,
		Name:   body.Name,
		Passwd: body.Passwd,
	}
	account.Token, _ = account.GenerateToken()

	account.HashPassword()
	if result := db.Create(&account); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error))
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   account.Token,
		Expires: time.Now().Add(7 * 24 * time.Hour),
		Path:    "/api",
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, api.Return("Created", echo.Map{
		"account":      account,
		"cookie_token": account.Token,
	}))
}

// @Summary check email's existense
// @Description
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Passwd path string true "user password"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/checkemail [POST]
func (h *AccountHandler) CheckEmail(c echo.Context) error {
	type RequestBody struct {
		Email string `json:"email" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	var account Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	} else {
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": true}))
	}
}

// @Summary login using email and passwd
// @Description
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Passwd path string true "user password"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/login [POST]
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
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(body.Passwd)) != nil {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
	}

	token, _ := account.GenerateToken()
	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(7 * 24 * time.Hour),
		Path:    "/api",
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
// @Router /account/{id}/logout [POST]
func (h *AccountHandler) LogoutAccount(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusBadRequest, api.Return("Not Logged in", nil))
	}
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	// cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	cookie.Path = "/api"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, api.Return("Account logged out", nil))
}

// @Summary the interface of modifying password
// @Description can only be called during logged-in status since there is no password check
// @Tags Account
// @Produce json
// @Param Email path string true "user e-mail"
// @Param Passwd path string true "user password (the new one)"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/{id}/modifypasswd [POST]
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
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
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
 * @brief public method for getting current logged-in account's ID.
 */
func getAccountID(c echo.Context) (string, error) {
	cookie, err := c.Cookie("token")
	if err != nil || cookie.Value == "" {
		return "", fmt.Errorf("Not logged in")
	}

	db, _ := c.Get("db").(*gorm.DB)
	var account Account
	if err := db.Where("token = ?", cookie.Value).First(&account).Error; err != nil { // not found
		return "", fmt.Errorf("Not logged in")
	}
	return account.ID, nil
}

/**
 * @brief middleware for getting current logged-in account's ID.
 */
func CheckAccountID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := getAccountID(c)
		if err != nil {
			return c.JSON(403, api.Return("unauthorised", err))
		}
		c.Set("id", id)
		return next(c)
	}
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
