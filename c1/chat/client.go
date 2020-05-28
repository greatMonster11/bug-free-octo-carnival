package main

import (
	"github.com/gorilla/websocket"
)

// client represent single chatting user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a chanel on which message are sent
	sent chan []byte
	// room is the room this client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.sent {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
