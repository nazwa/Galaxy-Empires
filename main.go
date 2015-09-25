package main

import (
	debug "bitbucket.org/tidepayments/gohelpers/gin"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
	"github.com/nazwa/galaxy-empires/ge"
	"github.com/nazwa/galaxy-empires/middleware"
)

var (
	ROOT_DIR string

	GE *ge.GalaxyEmpires
)

func main() {

	ROOT_DIR, _ = osext.ExecutableFolder()

	Universe = NewUniverseStruct(1, 15, 5)

	BaseData = NewBaseDataStore("data/buildings.json", "data/research.json")
	PlayerData = NewPlayerDataStore("data/players.json")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Errors("", "", nil))

	debug.AssignDebugHandlers(r.Group("/debug"))

	NewAccountHandler(r.Group("/account"))
	NewPlayerHandler(r.Group("/player", middleware.Authentication(JWTKey)))
	NewPlanetHandler(r.Group("/planet", middleware.Authentication(JWTKey)))

	r.Static("/assets", ROOT_DIR+"/web/assets")
	r.StaticFile("/", ROOT_DIR+"/web/index.html")

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
