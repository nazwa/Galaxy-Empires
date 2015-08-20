package main

import (
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	player := &PlayerStruct{
		Email:    r.PostFormValue("email"),
		Name:     r.PostFormValue("name"),
		Password: r.PostFormValue("password"),
	}

	if err := player.Validate(); err != nil {
		fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
		return
	}

	if err := player.HashPassword(); err != nil {
		fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
		return
	}
	PlayerData.SafeAdd(player)
	if token, err := player.CreateLoginToken(); err != nil {
		fmt.Fprintf(w, "{\"token\":\"%s\"}", token)
	} else {
		fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
	}

}
