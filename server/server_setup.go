package server

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) withAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sidCookie, err := r.Cookie(s.session.CookieName)
		if err != nil || sidCookie == nil {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}
		user, ok := s.session.Get(sidCookie.Value)
		if !ok {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), user, user)
		//create a new request using that new context
		rWithUser := r.WithContext(ctxWithUser)
		f(w, rWithUser)
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	user, err := s.userStorage.Login(login)
	if err != nil {
		log.Printf("User '%s' login failed\n", login)
		http.Error(w, "wrong login credentials", http.StatusUnauthorized)
		return
	}
	sid := s.session.Start(user.Id)
	cookie := &http.Cookie{
		Name:       s.session.CookieName,
		Value:      sid,
		Expires:    time.Now().Add(s.session.Ttl),
	}
	http.SetCookie(w, cookie)
	log.Printf("User '%s' logged in, id = %d\n", login, user.Id)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	user, err := s.userStorage.Register(login)
	if err != nil {
		http.Error(w, "wrong login credentials", http.StatusUnauthorized)
		return
	}
	sid := s.session.Start(user.Id)
	cookie := &http.Cookie{
		Name:       s.session.CookieName,
		Value:      sid,
		Expires:    time.Now().Add(s.session.Ttl),
	}
	http.SetCookie(w, cookie)
	log.Printf("User '%s' registered, id = %d\n", login, user.Id)
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client open ws connection")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
	}
	client := NewClient(uint32(id), ws, s)
	s.connectClient(client)
	client.Listen()
}

func (s *Server) Setup() {
	go s.ListenClients()
	http.HandleFunc("/ws", s.withAuth(s.wsHandler))
	http.HandleFunc("/login", s.loginHandler)
	http.HandleFunc("/register", s.registerHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
