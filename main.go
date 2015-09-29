package main

import (
	debug "bitbucket.org/tidepayments/gohelpers/gin"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
	"github.com/nazwa/galaxy-empires/config"
	"github.com/nazwa/galaxy-empires/ge"
	"github.com/nazwa/galaxy-empires/handlers"
	"github.com/nazwa/galaxy-empires/middleware"
)

var (
	ROOT_DIR string

	GE *ge.GalaxyEmpires
)

func main() {

	ROOT_DIR, _ = osext.ExecutableFolder()
	config.LoadConfig(ge.BuildFullPath(ROOT_DIR, "config.json"))

	GE = ge.NewGalaxyEmpires(ge.BuildFullPath(ROOT_DIR, "data"), ge.CoordinatesStruct{1, 15, 5})

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Errors("", "", nil))

	debug.AssignDebugHandlers(r.Group("/debug"))

	handlers.NewAccountHandler(r.Group("/account"), GE)
	handlers.NewPlayerHandler(r.Group("/player", middleware.Authentication([]byte(config.Config.Key))), GE)
	handlers.NewPlanetHandler(r.Group("/planet", middleware.Authentication([]byte(config.Config.Key))), GE)

	r.Static("/assets", ROOT_DIR+"/web/assets")
	r.StaticFile("/", ROOT_DIR+"/web/index.html")

	if err := r.Run(":" + config.Config.Port); err != nil {
		panic(err)
	}
}
