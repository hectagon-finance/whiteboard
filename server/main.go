package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"

	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/crypto"
)

type Input struct {
	PublicKey crypto.PublicKey
	SigNature *crypto.Signature
	Data      []byte
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	gob.Register(elliptic.P256())
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		buf := bytes.NewBuffer(message)
		dec := gob.NewDecoder(buf)

		var m Input

		if err := dec.Decode(&m); err != nil {
			log.Fatal(err)
		}

		log.Println(m.SigNature)
		log.Println(m.PublicKey)
		log.Println(m.Data)

		oke := m.SigNature.Verify(m.PublicKey, m.Data)
		log.Println(oke)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
