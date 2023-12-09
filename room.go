package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 512
)

type room struct {
	users map[*user]bool // Holds all current users in the room

	join chan *user // Add users to join in this room

	leave chan *user // user wishing to leave the channel

	forward chan []byte //Holds the incoming message to forwards to other client
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *user),
		leave:   make(chan *user),
		users:   make(map[*user]bool),
	}
}

func (r *room) run() {

	for {
		select {
		case user := <-r.join:
			r.users[user] = true
		case user := <-r.leave:
			delete(r.users, user)
			close(user.receive)
		case msg := <-r.forward:
			for user := range r.users {
				user.receive <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("serveHttp : ", err)
		return
	}

	user := &user{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- user
	defer func() { r.leave <- user }()

	go user.write()
	user.read()
}
