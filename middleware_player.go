package main

import (
	"bitbucket.org/nazwa/galaxy-empires/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	PlayerObjectKey string = "PlayerObject"
)

var (
	ErrorPlayerNotFound error = errors.New("Player not found!")
)

// This method checks if enpoints are logged in
func PlayerMiddleware(players *PlayerDataStore) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.MustGet(middleware.AuthUserIDKey).(string)
		if player := players.SafeGet(id); player == nil {
			c.AbortWithError(http.StatusInternalServerError, ErrorPlayerNotFound).SetType(gin.ErrorTypePublic)
			return
		} else {
			c.Set(PlayerObjectKey, player)
		}
	}

}
