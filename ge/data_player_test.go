package ge

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerStore(t *testing.T) {
	assert.Panics(t, func() {
		NewPlayerDataStore("")
	})
	assert.Panics(t, func() {
		NewPlayerDataStore("wrong_file.json")
	})

	var store *PlayerDataStore

	assert.NotPanics(t, func() {
		store = NewPlayerDataStore("../data/players.json")
	})
	assert.NotNil(t, store)
	assert.IsType(t, &PlayerDataStore{}, store)

}

func TestPlayerStoreOperations(t *testing.T) {
	store := NewPlayerDataStore("../data/players.json")

	assert.IsType(t, &PlayerDataStore{}, store)

	dummyPlayer := &PlayerStruct{}

	assert.Empty(t, store.Players)
	// Start by adding player with no id
	assert.Nil(t, store.SafeAdd(dummyPlayer))
	// We should have received an ID
	assert.NotNil(t, dummyPlayer.ID)
	// Fetching by that ID should return the same pointer
	assert.Equal(t, dummyPlayer, store.SafeGet(dummyPlayer.ID))
	// Trying to add this player again should fail
	assert.Error(t, store.SafeAdd(dummyPlayer))

	// Make sure getters can also fail nicely
	assert.Nil(t, store.SafeGet("invalid ID"))
	assert.Nil(t, store.SafeGetByEmail("adwaff sefsefsef"))

	// Check if getting by email works
	dummyPlayer.Email = "test@test.com"
	assert.Equal(t, dummyPlayer, store.SafeGetByEmail("test@test.com"))

	// Make sure removal works
	store.SafeRemove(dummyPlayer.ID)
	assert.Nil(t, store.SafeGet(dummyPlayer.ID))
	assert.Empty(t, store.Players)

	// Try adding with existing ID
	dummyPlayer2 := &PlayerStruct{
		ID:    "testID",
		Email: "hello@test.com",
	}
	dummyPlayer3 := &PlayerStruct{
		ID:    "testID",
		Email: "hello@test.com",
	}

	store.SafeAdd(dummyPlayer2)
	assert.Equal(t, dummyPlayer2, store.SafeGet(dummyPlayer2.ID))
	assert.Error(t, store.SafeAdd(dummyPlayer3))

}
