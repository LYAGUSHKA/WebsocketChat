package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Garius6/websocket_chat/pkg/server"
	"github.com/Garius6/websocket_chat/pkg/sockets"
	"github.com/Garius6/websocket_chat/pkg/storage/sqlstorage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
	SigningKey  string `json:"signing_key"`
	Port        string `json:"port"`
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

func configureChat(c *Config) (*server.Chat, error) {
	db, err := newDB(c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	store := sqlstorage.New(db)

	return &server.Chat{
		Store:        store,
		Rooms:        make(map[int]*sockets.Room),
		SessionStore: sessions.NewCookieStore([]byte(c.SigningKey)),
		SessionName:  "Auth",
		CtxKey:       0,
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

	r.HandleFunc("/", chat.LoginHandler)
	r.HandleFunc("/index", chat.AuthenticateUser(chat.ServeHome))
	r.HandleFunc("/user/create", chat.RegisterHandler).Methods("POST")
	r.HandleFunc("/ws/{chatID}", chat.RoomHandler)

	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+config.Port, handlers.LoggingHandler(os.Stdout, r)))
}
