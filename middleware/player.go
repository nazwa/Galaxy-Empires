package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"net/http"
)

// This method checks if enpoints are logged in
func PlayerMiddleware(players *ge.PlayerDataStore) gin.HandlerFunc {

	// Make sure we have players. It's OK to panic at this stage
	if players == nil {
		panic(ge.ErrorPlayerDatabaseMissing)
	}

	return func(c *gin.Context) {
		id := c.MustGet(AuthUserIDKey).(string)
		if player := players.SafeGet(id); player == nil {
			c.AbortWithError(http.StatusInternalServerError, ge.ErrorPlayerNotFound).SetType(gin.ErrorTypePublic)
			return
		} else {
			c.Set(ge.PlayerObjectKey, player)
		}
	}

}
