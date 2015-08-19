package main

import (
	"fmt"
	"github.com/kardianos/osext"
	"log"
	"net/http"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

var (
	ROOT_DIR string
	Store    *BaseDataStore
	JWTKey   []byte = []byte("This is a key")
)

func main() {
	ROOT_DIR, _ = osext.ExecutableFolder()

	Store = NewBaseDataStore("data/buildings.json", "data/research.json")
	for _, item := range Store.Buildings {
		fmt.Println(item)
	}

	http.Handle("/socket/", sockjs.NewHandler("/socket", sockjs.DefaultOptions, routerHandler))
	http.Handle("/", http.FileServer(http.Dir("web/")))

	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
