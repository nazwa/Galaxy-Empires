package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	ErrorInvalidCredentials error = errors.New("Invalid credentials")
)

type AccountHandler struct {
	RouterGroup *gin.RouterGroup
}

func (u *AccountHandler) Routes() {
	u.RouterGroup.POST("/login", gin.Bind(LoginStruct{}), u.LoginHandler)
	u.RouterGroup.POST("/register", gin.Bind(PlayerStruct{}), u.RegisterHandler)
}

func NewAccountHandler(r *gin.RouterGroup) *AccountHandler {
	u := &AccountHandler{RouterGroup: r}
	u.Routes()
	return u
}


func (u *AccountHandler) LoginHandler(c *gin.Context) {
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
	u.sendLoginToken(c, player)

}

func (u *AccountHandler) RegisterHandler(c *gin.Context) {
	player := c.MustGet(gin.BindKey).(*PlayerStruct)

	if err := player.HashPassword(); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	PlayerData.SafeAdd(player)
	u.sendLoginToken(c, player)
}

func (u *AccountHandler) sendLoginToken(c *gin.Context, player *PlayerStruct) {
	if token, err := player.CreateLoginToken(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"Token": token})
	}
}
