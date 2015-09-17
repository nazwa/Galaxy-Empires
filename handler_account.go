package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	ErrorInvalidCredentials error = errors.New("Invalid credentials")
)

func LoginHandler(c *gin.Context) {
	var player *PlayerStruct

	login := c.MustGet(gin.BindKey).(*LoginStruct)

	if player = PlayerData.SafeGetByEmail(login.Email); player == nil {
		c.AbortWithError(http.StatusUnauthorized, ErrorInvalidCredentials).SetType(gin.ErrorTypePublic)
		return
	}

	if err := player.CheckPassword(login.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.AbortWithError(http.StatusUnauthorized, ErrorInvalidCredentials).SetType(gin.ErrorTypePublic)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	sendLoginToken(c, player)

}

func RegisterHandler(c *gin.Context) {
	player := c.MustGet(gin.BindKey).(*PlayerStruct)

	fmt.Println(player)

	if err := player.HashPassword(); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	PlayerData.SafeAdd(player)
	sendLoginToken(c, player)
}

func sendLoginToken(c *gin.Context, player *PlayerStruct) {
	if token, err := player.CreateLoginToken(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"Token": token})
	}
}
