package main

import (
	"encoding/base64"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan *message

	//join is a channel for clients wishing to join the room
	join chan *client

	//leave is a channel for clients wishing to leave the room
	leave chan *client

	//clients holds all current clients in this room
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//joining
			r.clients[client] = true
		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					//send the message
				default:
					//failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("[Chat] ServeHTTP: ", err)
		return
	}

	cookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	user, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Fatal("Failed to get name from cookie:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: map[string]interface{}{"name": string(user)},
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
