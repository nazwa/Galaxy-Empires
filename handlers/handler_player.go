package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"github.com/nazwa/galaxy-empires/middleware"
	"net/http"
)

type PlayerHandler struct {
	RouterGroup *gin.RouterGroup
	GE          *ge.GalaxyEmpires
}

func (p *PlayerHandler) Routes() {
	p.RouterGroup.Use(middleware.PlayerMiddleware(p.GE.PlayerData))
	p.RouterGroup.GET("", p.Get)
	p.RouterGroup.GET("/data", p.GetData)
}

func NewPlayerHandler(r *gin.RouterGroup, ge *ge.GalaxyEmpires) *PlayerHandler {
	p := &PlayerHandler{
		RouterGroup: r,
		GE:          ge,
	}
	p.Routes()
	return p
}

func (p *PlayerHandler) Get(c *gin.Context) {
	player := c.MustGet(ge.PlayerObjectKey).(*ge.PlayerStruct)

	c.JSON(http.StatusOK, player.ToPublic(true))
}

func (p *PlayerHandler) GetData(c *gin.Context) {

	c.JSON(http.StatusOK, p.GE.BaseData)
}
