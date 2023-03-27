package validator

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

func ClientHandler(v *Validator, is_genesis string) {
	if is_genesis != "genesis" {
		u := url.URL{Scheme: "ws", Host: "localhost:" + is_genesis, Path: "/ws"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil) 
		defer conn.Close()
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
	}
	fmt.Println("Server is running on port: ", Port)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", v.Serve)
	log.Fatal(http.ListenAndServe("localhost:" + Port, nil))
}

func BroadcastBlockHash() {
	for {
		msg := <- Chan_1
		i := 0
		log.Println("** Client Handler** bool:", (DraftBlock.Hash == [32]byte{}) && ( (i == 10000) || (msg.memPool.Size() >= msg.memPool.CutOff)))
		if (DraftBlock.Hash == [32]byte{}) && ( (i == 10000) || (msg.memPool.Size() >= msg.memPool.CutOff)) {
			var k int

			if msg.memPool.Size() >= msg.memPool.CutOff {
				k = msg.memPool.CutOff
			} else {
				k = msg.memPool.Size()
			}

			DraftBlock = NewBlock(msg.heigt, [32]byte{}, msg.memPool.Transactions[:k])

			// broadcast block hash
			blockHashSlice := DraftBlock.Hash[:]
			blockHashStr := hex.EncodeToString(blockHashSlice)

			message := map[string]interface{}{
				"type":        "blockHash",
				"from": 	   Port,
				"blockHash":   blockHashStr,
			}

			Chan_2 <- k
			
			ConnectAndSendMessage(message)
		}
		time.Sleep(100 * time.Millisecond)
		i ++
	}
}