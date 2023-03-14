// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/crypto"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Input struct {
	PublicKey crypto.PublicKey
	SigNature *crypto.Signature
	Data      []byte
}

func main() {
	gob.Register(elliptic.P256())
	privKey := crypto.GeneratePrivateKey()
	publicKey := privKey.PublicKey()
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)

	input := Input{
		PublicKey: publicKey,
		SigNature: sig,
		Data:      msg,
	}

	fmt.Println(publicKey)
	fmt.Println(input)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(input); err != nil {
		return
	}

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	err = c.WriteMessage(websocket.TextMessage, buf.Bytes())
	if err != nil {
		log.Fatal("write:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
