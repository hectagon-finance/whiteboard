package validator

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func ClientHandler(v *Validator, is_genesis string) {
	if is_genesis != "genesis" {
		u := url.URL{Scheme: "ws", Host: "localhost:" + is_genesis, Path: "/ws"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil) 
		if err != nil {
			log.Fatal("dial:", err)
		}

		msg := map[string]interface{}{
			"type": "sync all request",
			"from": Port,
		}

		msgByte, err := json.Marshal(msg)
		if err != nil {
			log.Fatal("marshal:", err)
		}
		log.Println("Sending sync all request from", Port, "to", is_genesis)
		conn.WriteMessage(websocket.TextMessage, msgByte)
		conn.Close()
	}
	fmt.Println("Server is running on port: ", Port)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", v.Serve)
	log.Fatal(http.ListenAndServe("localhost:" + Port, nil))
}