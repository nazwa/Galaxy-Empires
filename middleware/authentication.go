package middleware

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	ErrorInvalidToken         = errors.New("Invalid or expired token")
	ErrorInvalidEncryptionKey = errors.New("Invalid encryption key")
)

const (
	AuthUserIDKey              string = "UserID"
	AuthTokenKey               string = "AuthToken"
	AuthTokenHeaderKey         string = "Token"
	AuthMinEncryptionKeyLength int    = 8
)

// This method checks if enpoints are logged in
func Authentication(encryptionKey []byte) gin.HandlerFunc {

	// This is a big problem that happens only on startup, let's panic
	if len(encryptionKey) < AuthMinEncryptionKeyLength {
		panic(ErrorInvalidEncryptionKey)
	}

	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get(AuthTokenHeaderKey)

		if len(tokenHeader) == 0 {
			c.AbortWithError(http.StatusUnauthorized, ErrorInvalidToken).SetType(gin.ErrorTypePublic)
			return
		}

		token, err := jwt.Parse(tokenHeader, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			if int64(t.Claims["exp"].(float64)) < time.Now().Unix() {
				return nil, ErrorInvalidToken
			}

			return encryptionKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, ErrorInvalidToken).SetType(gin.ErrorTypePublic)
			return
		}

		c.Set(AuthUserIDKey, token.Claims["id"])
		c.Set(AuthTokenKey, token)

	}
}
