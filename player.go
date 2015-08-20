package main

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type PlayerStruct struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (p *PlayerStruct) Validate() error {
	if l := len(p.Name); l < 1 || l > 50 {
		return errors.New("Invalid name")
	}
	if l := len(p.Email); l < 1 || l > 150 {
		return errors.New("Invalid email")
	}
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, p.Email); !m {
		return errors.New("Invalid email")
	}
	if l := len(p.Password); l < 5 || l > 50 {
		return errors.New("Invalid password")
	}
	return nil
}

func (p *PlayerStruct) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), 7)
	if err == nil {
		p.Password = string(hashedPassword)
	}
	return nil
}

func (p *PlayerStruct) CreateLoginToken() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = p.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(JWTKey)

}
