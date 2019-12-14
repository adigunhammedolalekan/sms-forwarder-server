package types

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"-"`
	Token string `json:"token" gorm:"-" sql:"-"`
}

type Token struct {
	jwt.StandardClaims
	UserId uint
	Email string
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}
