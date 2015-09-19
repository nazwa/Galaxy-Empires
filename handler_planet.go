package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	ErrorInvalidPlanetID error = errors.New("Invalid planet ID")
)

type PlanetHandler struct {
	RouterGroup *gin.RouterGroup
}

func (p *PlanetHandler) Routes() {
	p.RouterGroup.Use(PlayerMiddleware(PlayerData))
	p.RouterGroup.GET("/:id", p.Get)
}

func NewPlanetHandler(r *gin.RouterGroup) *PlanetHandler {
	p := &PlanetHandler{RouterGroup: r}
	p.Routes()
	return p
}

func (p *PlanetHandler) Get(c *gin.Context) {
	player := c.MustGet(PlayerObjectKey).(*PlayerStruct)

	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	planet := player.GetPlanet(id)
	if planet == nil {
		c.AbortWithError(http.StatusBadRequest, ErrorInvalidPlanetID).SetType(gin.ErrorTypePublic)
		return
	}

	planet.RecalculateResources(BaseData)
	c.JSON(http.StatusOK, planet)
}
