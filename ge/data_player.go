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

func (p *PlayerDataStore) SafeAdd(player *PlayerStruct) {
	p.m.Lock()
	defer p.m.Unlock()

	var token string

	// Find a new random id
	for {
		token = <-tokens.Token24
		if _, ok := p.Players[token]; !ok {
			break
		}
	}
	player.ID = token
	p.Players[token] = player
}

func (p *PlayerDataStore) SafeRemove(id string) {
	p.m.RLock()
	defer p.m.RUnlock()

	delete(p.Players, id)
}

func (p *PlayerDataStore) SafeGet(id string) *PlayerStruct {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.Players[id]
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
