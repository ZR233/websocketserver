// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocketserver

import (
	"net/http"
	"sync/atomic"
)

type Clients map[uint64]ClientInterface

// Center maintains the set of active Clients and broadcasts messages to the
// Clients.
type Center struct {
	clientIdIter uint64

	// Registered Clients.
	clients Clients

	// Inbound messages from the Clients.
	broadcast chan []byte

	// Register requests from the Clients.
	register chan ClientInterface

	// Unregister requests from Clients.
	unregister chan uint64

	msgHandler func(ClientInterface, []byte)
}

func (c *Center) Clients() Clients {
	return c.clients
}

func NewCenter(msgHandler func(ClientInterface, []byte)) *Center {
	return &Center{
		broadcast:  make(chan []byte),
		register:   make(chan ClientInterface),
		unregister: make(chan uint64),
		clients:    make(Clients),
		msgHandler: msgHandler,
	}
}

func (c *Center) registerClient(client ClientInterface) {
	c.clients[client.GetId()] = client
}

func (c *Center) unregisterClient(id uint64) {
	if v, ok := c.clients[id]; ok {
		v.Delete()
		delete(c.clients, id)
	}
}

func (c *Center) Run() {
	for {
		select {
		case client := <-c.register:
			c.registerClient(client)
		case id := <-c.unregister:
			c.unregisterClient(id)
		}
	}
}

// serveWs handles websocket requests from the peer.
func (c *Center) AddClient(client ClientInterface, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	sendChan := make(chan []byte, 256)
	clientId := atomic.AddUint64(&c.clientIdIter, 1)
	client = client.new(clientId, c, conn, &sendChan, c.msgHandler)
	//client := &Client{center: c, conn: conn, send: make(chan []byte, 256), msgHandler: c.msgHandler}
	client.run()
	c.register <- client
	return nil
}
