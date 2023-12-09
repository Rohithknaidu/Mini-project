package main

import (
	"github.com/gorilla/websocket"
)

type user struct {
	// Userspecific Websocket
	socket *websocket.Conn

	receive chan []byte // Channel to receive the messages

	room *room // user is chatting in which room
}

func (user *user) read() {

	defer user.socket.Close()

	// Keep on reading the messahe
	for {

		_, msg, err := user.socket.ReadMessage()
		if err != nil {
			return
		}
		user.room.forward <- msg // Read the messahe and write to the room channel
	}
}

func (user *user) write() {

	defer user.socket.Close()

	for msg := range user.receive {
		err := user.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
