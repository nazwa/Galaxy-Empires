package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"net/http"
	"strconv"
)

// This method checks if enpoints are logged in
func PlanetMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		player := c.MustGet(ge.PlayerObjectKey).(*ge.PlayerStruct)

		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
			return
		}

		planet := player.GetPlanet(id)
		if planet == nil {
			c.AbortWithError(http.StatusBadRequest, ge.ErrorInvalidPlanetID).SetType(gin.ErrorTypePublic)
			return
		}
		c.Set(ge.PlanetObjectKey, planet)
	}
}
