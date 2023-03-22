package start

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func ConnectAndSendMessage(v *ValidatorStruct, message map[string]interface{}) {
	v.ClientsMutex.Lock()
	defer v.ClientsMutex.Unlock()

	for _, peer := range v.Peers {
		if peer != strconv.Itoa(v.Port) {
			// Check if the peer is already connected
			isConnected := false
			for conn := range v.Clients {
				if conn.RemoteAddr().String() == "localhost:"+peer {
					isConnected = true
					break
				}
			}
			if isConnected {
				continue
			}

			u := url.URL{Scheme: "ws", Host: "localhost:" + peer, Path: "/ws"}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				fmt.Println("Error connecting to the peer:", err)
				continue
			}
			defer conn.Close()

			v.Clients[conn] = true

			sendMessage(v, conn, message)
		}
	}
}

func sendMessage(v *ValidatorStruct, conn *websocket.Conn, message map[string]interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling the message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}