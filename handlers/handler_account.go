package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/config"
	"github.com/nazwa/galaxy-empires/ge"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AccountHandler struct {
	RouterGroup *gin.RouterGroup
	GE          *ge.GalaxyEmpires
}

func (u *AccountHandler) Routes() {
	u.RouterGroup.POST("/login", gin.Bind(ge.LoginStruct{}), u.LoginHandler)
	u.RouterGroup.POST("/register", gin.Bind(ge.PlayerStruct{}), u.RegisterHandler)
}

func NewAccountHandler(r *gin.RouterGroup, ge *ge.GalaxyEmpires) *AccountHandler {
	u := &AccountHandler{
		RouterGroup: r,
		GE:          ge,
	}
	u.Routes()
	return u
}

func (u *AccountHandler) LoginHandler(c *gin.Context) {
	var player *ge.PlayerStruct

	login := c.MustGet(gin.BindKey).(*ge.LoginStruct)

	if player = u.GE.PlayerData.SafeGetByEmail(login.Email); player == nil {
		c.AbortWithError(http.StatusUnauthorized, ge.ErrorInvalidCredentials).SetType(gin.ErrorTypePublic)
		return
	}

	if err := player.CheckPassword(login.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.AbortWithError(http.StatusUnauthorized, ge.ErrorInvalidCredentials).SetType(gin.ErrorTypePublic)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	u.sendLoginToken(c, player)

}

func (u *AccountHandler) RegisterHandler(c *gin.Context) {
	player := c.MustGet(gin.BindKey).(*ge.PlayerStruct)

	if err := player.HashPassword(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Make sure this email hasn't been used yet
	if u.GE.PlayerData.SafeGetByEmail(player.Email) != nil {
		c.AbortWithError(http.StatusBadRequest, ge.ErrorEmailInUse).SetType(gin.ErrorTypePublic)
		return
	}

	u.GE.PlayerData.SafeAdd(player)

	// We have a player, let's make him a planet!
	planet, err := ge.GenerateNewPlanet(u.GE.Universe, u.GE.BaseData)
	if err != nil {
		u.GE.PlayerData.SafeRemove(player.ID)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	player.AddPlanet(planet)

	u.sendLoginToken(c, player)
}

func (u *AccountHandler) createLoginToken(player *ge.PlayerStruct) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = player.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(config.Config.Key))

}

func (u *AccountHandler) sendLoginToken(c *gin.Context, player *ge.PlayerStruct) {
	if token, err := u.createLoginToken(player); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"Token": token})
	}
}
