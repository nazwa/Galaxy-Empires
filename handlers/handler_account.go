package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Make sure this email hasn't been used yet
	if PlayerData.SafeGetByEmail(player.Email) != nil {
		c.AbortWithError(http.StatusBadRequest, ErrorEmailInUse).SetType(gin.ErrorTypePublic)
		return
	}

	PlayerData.SafeAdd(player)

	// We have a player, let's make him a planet!
	planet, err := GenerateNewPlanet(Universe, BaseData)
	if err != nil {
		PlayerData.SafeRemove(player.ID)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	player.AddPlanet(planet)

	u.sendLoginToken(c, player)
}

func (u *AccountHandler) createLoginToken(player *ge.PlayerStruct) string {

	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = player.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(JWTKey)

}

func (u *AccountHandler) sendLoginToken(c *gin.Context, player *PlayerStruct) {
	if token, err := u.createLoginToken(player); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"Token": token})
	}
}
