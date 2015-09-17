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
	JWTKey     []byte = []byte("This is a key")
)

func main() {
	ROOT_DIR, _ = osext.ExecutableFolder()

	BaseData = NewBaseDataStore("data/buildings.json", "data/research.json")
	PlayerData = NewPlayerDataStore("data/players.json")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.Errors("", "", nil))
	r.Use(gin.Recovery())

	r.POST("/account/login", gin.Bind(LoginStruct{}), LoginHandler)
	r.POST("/account/register", gin.Bind(PlayerStruct{}), RegisterHandler)

	r.Static("/", ROOT_DIR+"/web")

	log.Println("Server started on port: 8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
