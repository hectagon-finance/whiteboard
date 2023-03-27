package validator

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func ConnectAndSendMessage(v *Validator, message map[string]interface{}) {
	for _, peer := range v.Peers {
		if peer != v.Port {
			u := url.URL{Scheme: "ws", Host: "localhost:" + peer, Path: "/ws"}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				fmt.Println("Error connecting to the peer:", err)
				continue
			}
			
			msg, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error marshaling the message:", err)
					return
				}
			conn.WriteMessage(websocket.TextMessage, msg)

			conn.Close()
		}
	}
}
