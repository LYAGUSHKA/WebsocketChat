package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Garius6/websocket_chat/storage"
	"github.com/gorilla/mux"
)

var db *storage.Storage
var mySigningKey = []byte("secret")

func init() {
	db := storage.New(&storage.Config{"message.sqlite3"})
	_ = db
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/templates/chat.html")
}

func main() {
	rooms := make(map[int]*Room)

	r := mux.NewRouter()

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws/{chatID}", func(w http.ResponseWriter, r *http.Request) {

		//Parse route parameter
		vars := mux.Vars(r)
		chatID, err := strconv.Atoi(vars["chatID"])
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		var room *Room
		//If route correct and chat exist -> connecting to chat
		//else create new chat
		if _, ok := rooms[chatID]; ok {
			room = rooms[chatID]
		} else {
			rooms[chatID] = newRoom(db)
			room = rooms[chatID]
			go rooms[chatID].run()
			for _, v := range rooms {
				v.events <- Event{NEWCHAT, struct {
					ID   int
					chat *Room
				}{chatID, room}}
			}
		}

		serveWs(room, w, r)
	})

	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
