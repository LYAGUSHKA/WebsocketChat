package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/storage"
	"github.com/gorilla/mux"
)

type Chat struct {
	Store      *storage.Storage
	Rooms      map[int]*Room
	SigningKey []byte
}

func New() *Chat {
	return &Chat{}
}

func (c *Chat) serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/templates/chat.html")
}

func (c *Chat) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/templates/auth.html")
	} else if r.Method == "POST" {
		nickname := r.FormValue("login")
		password := r.FormValue("password")
		c.Store.User().Create(&model.User{Nickname: nickname, EncryptedPassword: password})
	}
}

func (c *Chat) roomHandler(w http.ResponseWriter, r *http.Request) {

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
	if _, ok := c.Rooms[chatID]; ok {
		room = c.Rooms[chatID]
	} else {
		c.Rooms[chatID] = newRoom(c.Store)
		room = c.Rooms[chatID]
		go c.Rooms[chatID].run()
		for _, v := range c.Rooms {
			v.events <- Event{NEWCHAT, struct {
				ID   int
				chat *Room
			}{chatID, room}}
		}
	}

	serveWs(room, w, r)
}

func main() {
	c := New()
	c.Rooms = make(map[int]*Room)
	c.Store = storage.New(&storage.Config{DatabaseURL: "message.sqlite3"})
	if err := c.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/", c.serveHome)
	r.HandleFunc("/register", c.registerHandler)
	r.HandleFunc("/ws/{chatID}", c.roomHandler)

	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
