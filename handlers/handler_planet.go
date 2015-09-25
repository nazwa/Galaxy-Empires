package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"github.com/nazwa/galaxy-empires/middleware"
	"net/http"
	"time"
)

type PlanetHandler struct {
	RouterGroup *gin.RouterGroup
	GE          *ge.GalaxyEmpires
}

type PlanetRenameForm struct {
	Name string `form:"name" binding:"required,min=1,max=60"`
}

func (p *PlanetHandler) Routes() {
	p.RouterGroup.Use(middleware.PlayerMiddleware(p.GE.PlayerData))
	p.RouterGroup.Use(middleware.PlanetMiddleware())

	p.RouterGroup.GET("/:id", p.Get)
	p.RouterGroup.POST("/:id/rename", gin.Bind(PlanetRenameForm{}), p.Rename)
	p.RouterGroup.POST("/:id/building/build/:type", p.BuildBuilding)
	p.RouterGroup.POST("/:id/building/cancel", p.CancelBuilding)
}

func NewPlanetHandler(r *gin.RouterGroup, ge *ge.GalaxyEmpires) *PlanetHandler {
	p := &PlanetHandler{
		RouterGroup: r,
		GE:          ge,
	}
	p.Routes()
	return p
}

func (p *PlanetHandler) Get(c *gin.Context) {
	planet := c.MustGet(ge.PlanetObjectKey).(*ge.PlanetStruct)

	planet.UpdatePlanet(p.GE.BaseData, time.Now())
	c.JSON(http.StatusOK, planet.ToPublic(true))
}

func (p *PlanetHandler) Rename(c *gin.Context) {
	planet := c.MustGet(ge.PlanetObjectKey).(*ge.PlanetStruct)
	rename := c.MustGet(gin.BindKey).(*PlanetRenameForm)
	planet.Name = rename.Name

	c.JSON(http.StatusOK, planet.ToPublic(true))
}

func (p *PlanetHandler) BuildBuilding(c *gin.Context) {
	planet := c.MustGet(ge.PlanetObjectKey).(*ge.PlanetStruct)
	building_id := c.Params.ByName("type")

	if err := planet.BuildBuilding(p.GE.BaseData, building_id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, ge.DefaultSuccessResponse)
}

func (p *PlanetHandler) CancelBuilding(c *gin.Context) {
	planet := c.MustGet(ge.PlanetObjectKey).(*ge.PlanetStruct)
	planet.CancelBuilding()

	c.JSON(http.StatusOK, ge.DefaultSuccessResponse)
}
