package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/middleware"
	"net/http"
)

// This method checks if enpoints are logged in
func PlayerMiddleware(players *PlayerDataStore) gin.HandlerFunc {

	// Make sure we have players. It's OK to panic at this stage
	if players == nil {
		panic(ErrorPlayerDatabaseMissing)
	}

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
