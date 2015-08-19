package main

import (
	"encoding/json"
)

type SocketContext struct {
	Command string          `json:"command"`
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
	Session sockjs.Session
}

var (
	ErrorInternalServerError error = errors.New("{\"error\":\"Woops! Something went wrong.\"}")
)

func (c *SocketContext) JSON(v interface{}, session sockjs.Session) error {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	return session.Send(string(data))
}

func InternalServerError(err error, session sockjs.Session) {
	// LOG ERROR
	log.Println(err)

	session.Send(ErrorInternalServerError.Error())
}
