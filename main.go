package main

import (
	"bitbucket.org/nazwa/galaxy-empires/middleware"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
	"log"
)

var (
	ROOT_DIR   string
	BaseData   *BaseDataStore
	PlayerData *PlayerDataStore
	Universe *UniverseStruct
	JWTKey     []byte = []byte("This is a key")
)

func main() {
	Universe = NewUniverseStruct(1, 15, 5)
		
	
	ROOT_DIR, _ = osext.ExecutableFolder()

	BaseData = NewBaseDataStore("data/buildings.json", "data/research.json")
	PlayerData = NewPlayerDataStore("data/players.json")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Errors("", "", nil))

	NewAccountHandler(r.Group("/account"))
	NewPlayerHandler(r.Group("/player", middleware.Authentication(JWTKey)))

	r.Static("/assets", ROOT_DIR+"/web/assets")
	r.StaticFile("/", ROOT_DIR + "/web/index.html")

	log.Println("Server started on port: 8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
