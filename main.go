package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/home.html")
}

func main() {
	rooms := make(map[int]*Room)

	r := mux.NewRouter()

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws/{chatID}", func(w http.ResponseWriter, r *http.Request) {
		//Connect to database
		db, err := sql.Open("sqlite3", "message.sqlite3")
		if err != nil {
			_ = fmt.Errorf("%s", err)
			http.Error(w, "Something goes wrong", http.StatusInternalServerError)
			return
		}

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
		}

		serveWs(room, w, r)
	})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
