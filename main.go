package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName            = "gopher"
	ctxKey      ContextKey = iota
)

type ContextKey int8

type Chat struct {
	Store        *storage.Storage
	Rooms        map[int]*Room
	SigningKey   []byte
	SessionStore sessions.Store
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (c *Chat) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/templates/reg.html")
	} else if r.Method == "POST" {
		nickname := r.FormValue("login")
		password := r.FormValue("password")

		u, err := c.Store.User().FindByNickname(nickname)
		if err != nil {
			http.Error(w, "db", http.StatusInternalServerError)
			return
		}

		if password != u.EncryptedPassword {
			http.Error(w, "Password", http.StatusInternalServerError)
			return
		}

		session, err := c.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, "Coockie", http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = u.ID
		err = c.SessionStore.Save(r, w, session)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (c *Chat) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := c.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u, err := c.Store.User().FindByID(int(id.(uint64)))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey, u)))
	})
}

func main() {
	c := &Chat{}
	c.Rooms = make(map[int]*Room)
	c.SigningKey = []byte("secret")
	c.SessionStore = sessions.NewCookieStore(c.SigningKey)
	c.Store = storage.New(&storage.Config{DatabaseURL: "message.sqlite3"})
	if err := c.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/", c.authenticateUser(c.serveHome))
	r.HandleFunc("/login", c.loginHandler)
	r.HandleFunc("/user/create", c.registerHandler)
	r.HandleFunc("/ws/{chatID}", c.roomHandler)

	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
