package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	BCRYPT_COMPLEXITY int = 12
)

type PlayerStruct struct {
	ID       string `binding:"omitempty,number"`
	Name     string `form:"name" binding:"required,min=1,max=60"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=5,max=60"`
}

type LoginStruct struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"max=200"`
}

func (p *PlayerStruct) GenerateHash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), BCRYPT_COMPLEXITY)
	return string(hash), err
}
func (p *PlayerStruct) HashPassword() error {
	hashedPassword, err := p.GenerateHash(p.Password)
	if err == nil {
		p.Password = hashedPassword
	}
	return err
}

func (p *PlayerStruct) CreateLoginToken() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = p.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(JWTKey)

}

func (p *PlayerStruct) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
}
