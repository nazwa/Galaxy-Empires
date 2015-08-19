package main

import (
	"errors"
)

type CommandHandlerFunc func(c *SocketContext) error

var (
	handlers             map[string]CommandHandlerFunc
	ErrorCommandNotFound error = errors.New("Route not found")
)

func AddHandler(command string, fn CommandHandlerFunc) {
	if handlers == nil {
		handlers = make(map[string]CommandHandlerFunc)
	}
	handlers[command] = fn
}

func ExecuteCommand(c *SocketContext) error {
	if fn, ok := handlers[c.Command]; ok {
		return fn(c)
	} else {
		return ErrorCommandNotFound
	}
}
