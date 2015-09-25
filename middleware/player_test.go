package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nazwa/galaxy-empires/ge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerMiddleware(t *testing.T) {
	// Should panic if no store
	assert.Panics(t, func() {
		PlayerMiddleware(nil)
	})

	// This shouldnt
	store := ge.NewPlayerDataStore("../data/players.json")
	assert.NotPanics(t, func() {
		PlayerMiddleware(store)
	})

	handler := PlayerMiddleware(store)
	c := &gin.Context{}

	// First without a player id
	// Should panic. This assumes handler tried to add an error to writer
	assert.Panics(t, func() {
		handler(c)
	})

	// Then with ID but without the right player

	c.Set(AuthUserIDKey, "random key")
	assert.Panics(t, func() {
		handler(c)
	})

	player := &ge.PlayerStruct{}
	store.SafeAdd(player)

	// Now with a correct player ID
	c.Set(AuthUserIDKey, player.ID)
	assert.NotPanics(t, func() {
		handler(c)
	})
	assert.NotPanics(t, func() {
		assert.Equal(t, player, c.MustGet(ge.PlayerObjectKey))
	})

}
