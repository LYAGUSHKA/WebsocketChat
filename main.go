package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Garius6/websocket_chat/pkg/sockets"
	"github.com/Garius6/websocket_chat/pkg/storage"
	"github.com/Garius6/websocket_chat/pkg/storage/sqlstorage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName            = "Auth"
	ctxKey      ContextKey = iota
)

type ContextKey int8

type Config struct {
	DatabaseURL string `json:"database_url"`
	SigningKey  string `json:"signing_key"`
	Port        string `json:"port"`
}

type Chat struct {
	Store        storage.Storage
	Rooms        map[int]*sockets.Room
	SessionStore sessions.Store
}

func (c *Chat) serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/templates/chat.html")
}

func (c *Chat) registerHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	_, err := c.Store.User().Create(login, password)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (c *Chat) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/templates/reg.html")
	} else if r.Method == "POST" {
		nickname := r.FormValue("login")
		password := r.FormValue("password")

		u, err := c.Store.User().FindByLogin(nickname)
		if err != nil {
			http.Error(w, "db", http.StatusInternalServerError)
			return
		}

		if u.ComparePassword(password) {
			http.Error(w, "Password", http.StatusBadRequest)
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
	var room *sockets.Room
	//If route correct and chat exist -> connecting to chat
	//else create new chat
	if _, ok := c.Rooms[chatID]; ok {
		room = c.Rooms[chatID]
	} else {
		c.Rooms[chatID] = sockets.NewRoom(c.Store)
		room = c.Rooms[chatID]
		go c.Rooms[chatID].RunRoom()
		for _, v := range c.Rooms {
			v.Events <- sockets.Event{
				Type: sockets.NewChat,
				Data: sockets.ChatInfo{ID: chatID, Chat: room},
			}
		}
	}

	sockets.ServeWs(room, w, r)
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

func getConfig(configName string) (*Config, error) {
	file, err := os.Open(configName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c := Config{}
	if err = json.NewDecoder(file).Decode(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func configureChat(c *Config) (*Chat, error) {
	db, err := newDB(c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	store := sqlstorage.New(db)

	return &Chat{
		Rooms:        make(map[int]*sockets.Room),
		SessionStore: sessions.NewCookieStore([]byte(c.SigningKey)),
		Store:        store,
	}, nil
}

func main() {
	config, err := getConfig("configs/chat.json")
	if err != nil {
		log.Fatal(err)
	}

	chat, err := configureChat(config)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/", chat.authenticateUser(chat.serveHome))
	r.HandleFunc("/login", chat.loginHandler)
	r.HandleFunc("/user/create", chat.registerHandler).Methods("POST")
	r.HandleFunc("/ws/{chatID}", chat.roomHandler)

	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+config.Port, handlers.LoggingHandler(os.Stdout, r)))
}
