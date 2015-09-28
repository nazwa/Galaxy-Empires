package ge

import (
	"bitbucket.org/tidepayments/gohelpers/tokens"
	"strings"
	"sync"
)

type PlayerDataStore struct {
	Players map[string]*PlayerStruct
	m       sync.RWMutex
}

func NewPlayerDataStore(file string) *PlayerDataStore {
	store := &PlayerDataStore{
		Players: make(map[string]*PlayerStruct),
	}
	if err := LoadFile(file, &store.Players); err != nil {
		panic(err)
	}

	return store
}

// Add generated the ID if necesary.
// Email is not used, not sure if thats a good approahc or not.
// Maybe email should be the ID? That would solve two problems in one go
// But then referencing parent b email everywhere is a bit silly
func (p *PlayerDataStore) SafeAdd(player *PlayerStruct) error {
	p.m.Lock()
	defer p.m.Unlock()

	// If player already has an ID, use it. This might be handy on data load
	if len(player.ID) != 0 {
		// Make sure that ID doesn't already exist
		if _, ok := p.Players[player.ID]; ok {
			return ErrorPlayerIDInUse
		}
	} else {
		var token string
		// Generate a new random id
		for {
			token = <-tokens.Token24
			if _, ok := p.Players[token]; !ok {
				break
			}
		}
		player.ID = token
	}
	p.Players[player.ID] = player
	return nil
}

// @todo: This needs a lot of work. Should planets also be removed?
// What about all other lined assets?
// Is this the place to do it or will there be a command or something handling it?
func (p *PlayerDataStore) SafeRemove(id string) {
	p.m.RLock()
	defer p.m.RUnlock()

	delete(p.Players, id)
}

func (p *PlayerDataStore) SafeGet(id string) *PlayerStruct {
	p.m.RLock()
	defer p.m.RUnlock()

	if player, ok := p.Players[id]; !ok {
		return nil
	} else {
		return player
	}
}

func (p *PlayerDataStore) SafeGetByEmail(email string) *PlayerStruct {
	p.m.RLock()
	defer p.m.RUnlock()

	for _, player := range p.Players {
		if strings.EqualFold(email, player.Email) {
			return player
		}
	}
	return nil
}
