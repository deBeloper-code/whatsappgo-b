package main

import "github.com/gorilla/websocket"

// Represents a single chatting user
type client struct {
	socket *websocket.Conn
	// is a channel on which messages are sent.
	send chan []byte
	// is a room this client is chatting in.
	room *room
}

func (c *client) read() {
	// After everything socket will be closed.
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
	// After everything socket will be closed.
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
