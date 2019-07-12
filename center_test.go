// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocketserver

import (
	"log"
	"net/http"
	"strconv"
	"testing"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

type MyClient struct {
	Client
}

func TestGetNewServer(t *testing.T) {
	handle := func(c ClientInterface, msg []byte) {
		t.Log(strconv.FormatUint(c.GetId(), 10) + ":" + string(msg))
	}
	center := NewCenter(handle)

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		center.AddClient(MyClient{}, w, r)
	})

	go center.Run()
	err := http.ListenAndServe("localhost:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
