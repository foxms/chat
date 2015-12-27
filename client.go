package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// client represents a single chatting user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room this client is chatting in
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			log.Println("Read from client err:", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			log.Println("Write to client err:", err)
			break
		}
	}
	c.socket.Close()
}
