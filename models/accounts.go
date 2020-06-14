package models

import (
	u "backend/utils"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Token JWT claims struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account struct to hold user account data
type Account struct {
	AccountID uint   `gorm:"column:accountid; auto_increment"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token" gorm:"-"`
}

//Validate function validates account data
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &Account{}

	err := GetDB().Table("account").Where("email = ?", account.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create function creates account and stores it in the database
func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}
	//encrypts password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Table("account").Create(account)

	if account.AccountID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	tk := &Token{UserID: account.AccountID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = ""

	fmt.Printf("Account created in with the email address %s", account.Email) //server logs
	response := u.Message(true, "Account created successfully")
	response["account"] = account

	return response
}

//Login function checks details with account data in the database
func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("account").Where("email=?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	// login was correct
	account.Password = ""

	tk := &Token{UserID: account.AccountID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	fmt.Printf("Account logged in with the email address %s", account.Email)
	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetAccount function gets account data from the database using an accountID
func GetAccount(id uint) *Account {
	acc := &Account{}
	GetDB().Table("account").Where("accountID = ?", id).First(acc)
	if acc.Email == "" {
		return nil //user not found
	}

	acc.Password = ""
	return acc
}
