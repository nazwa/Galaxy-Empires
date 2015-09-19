package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlayerHandler struct {
	RouterGroup *gin.RouterGroup
}

func (p *PlayerHandler) Routes() {
	p.RouterGroup.Use(PlayerMiddleware(PlayerData))
	p.RouterGroup.GET("", p.Get)
	p.RouterGroup.GET("/data", p.GetData)
}

func NewPlayerHandler(r *gin.RouterGroup) *PlayerHandler {
	p := &PlayerHandler{RouterGroup: r}
	p.Routes()
	return p
}

func (p *PlayerHandler) Get(c *gin.Context) {
	player := c.MustGet(PlayerObjectKey).(*PlayerStruct)

	c.JSON(http.StatusOK, player)
}

func (p *PlayerHandler) GetData(c *gin.Context) {

	c.JSON(http.StatusOK, BaseData)
}
