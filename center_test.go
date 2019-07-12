// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocketserver

import (
	"reflect"
	"testing"
)

func TestGetNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *Center
	}{
		{"", &Center{broadcast: make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
