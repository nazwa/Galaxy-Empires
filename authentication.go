package main

import (
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

var (
	ErrorInvalidToken = errors.New("Invalid or expired token")
)

func CheckToken(c *SocketContext) error {
	if len(c.Token) == 0 {
		return ErrorInvalidToken
	}

	token, err := jwt.Parse(c.Token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		if int64(t.Claims["exp"].(float64)) < time.Now().Unix() {
			return nil, ErrorInvalidToken
		}

		return JWTKey, nil
	})

	if err != nil || !token.Valid {
		return ErrorInvalidToken
	}
	return nil
}
