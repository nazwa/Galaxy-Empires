package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetMiddleware(t *testing.T) {
	// Should not panic on call
	assert.NotPanics(t, func() {
		PlanetMiddleware()
	})

	c := &gin.Context{}
	handler := PlanetMiddleware()

	// This one panics because there is no player
	assert.Panics(t, func() {
		handler(c)
	})

	player := &PlayerStruct{ID: "p1"}
	c.Set(PlayerObjectKey, player)
	// Now there is a player, but no planets
	assert.Panics(t, func() {
		handler(c)
	})

	// Now there is a player, but no planet id
	planet := &PlanetStruct{}
	player.AddPlanet(planet)
	assert.Panics(t, func() {
		handler(c)
	})

	// This should work. Player, planet and a valid index
	c.Params = gin.Params{gin.Param{"id", "0"}}
	assert.NotPanics(t, func() {
		handler(c)
	})

	// Woops! wrong index
	c.Params = gin.Params{gin.Param{"id", "10"}}
	assert.Panics(t, func() {
		handler(c)
	})

	// Woops! wrong index
	c.Params = gin.Params{gin.Param{"id", "text"}}
	assert.Panics(t, func() {
		handler(c)
	})

}
