// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocketserver

import (
	"net/http"
)

type Clients map[*Client]bool

// Center maintains the set of active Clients and broadcasts messages to the
// Clients.
type Center struct {
	// Registered Clients.
	clients Clients

	// Inbound messages from the Clients.
	broadcast chan []byte

	// Register requests from the Clients.
	register chan *Client

	// Unregister requests from Clients.
	unregister chan *Client

	msgHandler func(*Client, []byte)
}

func (c *Center) Clients() Clients {
	return c.clients
}

func newCenter(msgHandler func(*Client, []byte)) *Center {
	return &Center{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		msgHandler: msgHandler,
	}
}

func (c *Center) Run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
			}
		case message := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (c *Center) NewClient(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{center: c, conn: conn, send: make(chan []byte, 256), msgHandler: c.msgHandler}
	client.center.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}
