package main

import (
	"github.com/kardianos/osext"
	"log"
	"net/http"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

var (
	ROOT_DIR   string
	BaseData   *BaseDataStore
	PlayerData *PlayerDataStore
	JWTKey     []byte = []byte("This is a key")
)

func main() {
	ROOT_DIR, _ = osext.ExecutableFolder()

	BaseData = NewBaseDataStore("data/buildings.json", "data/research.json")
	PlayerData = NewPlayerDataStore("data/players.json")

	http.HandleFunc("/account/login", LoginHandler)
	http.HandleFunc("/account/register", RegisterHandler)
	http.Handle("/socket/", sockjs.NewHandler("/socket", sockjs.DefaultOptions, routerHandler))
	http.Handle("/", http.FileServer(http.Dir("web/")))

	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
