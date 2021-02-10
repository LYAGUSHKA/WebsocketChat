package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Garius6/websocket_chat/pkg/sockets"
	"github.com/Garius6/websocket_chat/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Chat struct {
	Store        storage.Storage
	Rooms        map[int]*sockets.Room
	SessionStore sessions.Store
	SessionName  string
	CtxKey       int8
}

func (c *Chat) ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/templates/chat.html")
}

func (c *Chat) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	_, err := c.Store.User().Create(login, password)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (c *Chat) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/templates/reg.html")
	} else if r.Method == "POST" {
		nickname := r.FormValue("login")
		password := r.FormValue("password")

		u, err := c.Store.User().FindByLogin(nickname)
		if err != nil {
			log.Println("loginHandler: ", err)
			http.Error(w, "db", http.StatusInternalServerError)
			return
		}

		if u.ComparePassword(password) {
			http.Error(w, "Password", http.StatusBadRequest)
			return
		}

		session, err := c.SessionStore.Get(r, c.SessionName)
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
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
}

func (c *Chat) RoomHandler(w http.ResponseWriter, r *http.Request) {

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

func (c *Chat) AuthenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := c.SessionStore.Get(r, c.SessionName)
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

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), c.CtxKey, u)))
	})
}
