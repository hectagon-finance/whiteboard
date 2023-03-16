// // Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var (
	MessageString = 1
	MessagePeers  = 2
	MessageInput = 3
)

type Input struct {
	Id 		  string `json:"id"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data 	  []byte `json:"data"`
}

type Message struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Content []interface{} `json:"content"`
	Type    int		`json:"type"`
}

func StartClient() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: ":" + "8080", Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer c.Close()

	done := make(chan struct{})
	
	go func() {
		defer close(done)
		for {
			// _, message, err := c.ReadMessage()
			// if err != nil {
			// 	return
			// }
			// log.Printf("Client: %s", message)
		}
	}()

	// go n.ConnectToNode()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		return
		// 	}
		// }
		case <-ticker.C:
			Broadcast()	
		}
	}
}

// User client to connect to another node
func Broadcast(args ...interface{}) {
	// connect to the peer's server and send a message
	u := url.URL{Scheme: "ws", Host: "localhost" + ":" + "8080", Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Println("dial:", err)
		return
	}
	defer c.Close()
	
	// convert to string
	Id := strconv.Itoa(rand.Intn(5))
	
	input := Input{Id: Id, PublicKey: "publicKey", Signature: "signature", Data: []byte("data")}
	
	content := []interface{}{input}

	message := Message{From: "client", To: "8080", Content: content, Type: 3}
	c.WriteJSON(message)
}


func main(){
	StartClient()
}