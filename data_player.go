package main

import (
	"bitbucket.org/tidepayments/gohelpers/tokens"
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

func (p *PlayerDataStore) SafeGet(id string) *PlayerStruct {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.Players[id]
}