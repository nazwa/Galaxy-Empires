package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlanetHandler struct {
	RouterGroup *gin.RouterGroup
}

type PlanetRenameForm struct {
	Name string `form:"name" binding:"required,min=1,max=60"`
}

func (p *PlanetHandler) Routes() {
	p.RouterGroup.Use(PlayerMiddleware(PlayerData))
	p.RouterGroup.GET("/:id", PlanetMiddleware(), p.Get)
	p.RouterGroup.POST("/:id/rename", PlanetMiddleware(), gin.Bind(PlanetRenameForm{}), p.Rename)
	p.RouterGroup.POST("/:id/building/build/:type", PlanetMiddleware(), p.BuildBuilding)
}

func NewPlanetHandler(r *gin.RouterGroup) *PlanetHandler {
	p := &PlanetHandler{RouterGroup: r}
	p.Routes()
	return p
}

func (p *PlanetHandler) Get(c *gin.Context) {
	planet := c.MustGet(PlanetObjectKey).(*PlanetStruct)

	planet.RecalculateResources(BaseData)
	c.JSON(http.StatusOK, planet.ToPublic(true))
}

func (p *PlanetHandler) Rename(c *gin.Context) {
	planet := c.MustGet(PlanetObjectKey).(*PlanetStruct)
	rename := c.MustGet(gin.BindKey).(*PlanetRenameForm)
	planet.Name = rename.Name

	c.JSON(http.StatusOK, planet.ToPublic(true))
}

func (p *PlanetHandler) BuildBuilding(c *gin.Context) {
	//planet := c.MustGet(PlanetObjectKey).(*PlanetStruct)
	//building_id := c.Params.ByName("type")

}
