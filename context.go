package main

import (
	"encoding/json"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"errors"
	"log"
	
)

type SocketContext struct {
	Command string          `json:"command"`
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
	Session sockjs.Session
}

type H map[string]interface{}
var (
	ErrorInternalServerError error = errors.New("{\"error\":\"Woops! Something went wrong.\"}")
)

func (c *SocketContext) JSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Session.Send(string(data))
}

func (c *SocketContext) InternalServerError(err error) {
	log.Println(err)
	
	c.JSON(H{"error": "Woops! Something went wrong"})
}