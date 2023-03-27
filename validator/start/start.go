package start

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func Start(v *Validator) {
	go StartServer(v)

	go StartClient(v)
}

func Serve(w http.ResponseWriter, r *http.Request) {
	// Create a new Gorilla WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	
}

func ClientHandler(v *Validator){

}

func handleConnections(v *Validator, conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		HandleMessage(v, msg)
	}
}

func StartClient(v *Validator) {
	u := url.URL{Scheme: "ws", Host: "localhost:" + strconv.Itoa(v.Port), Path: "/ws"}

	// Retry connecting to the server with a delay
	for {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Error connecting to the server, retrying in 1 second:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		defer conn.Close()
		log.Println("Connected to the server")

		// Check if need to sync with the network
		// Sync(v)
		// time.Sleep(3 * time.Second)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		counterr := 0
		for {
			select {
			case <-ticker.C:
				// BroadcastMempool(v)
				BroadcastPeer(v)
				counterr = checkMemPool(v, counterr)		
			}
		}
	}
}

func checkMemPool(v *Validator, counter int) int {
	if v.MemPool.Size() == 3  || counter == 10 {
		v.TempBlock = NewBlock(1, [32]byte{}, v.MemPool.GetTransactions())
		blockHash := v.TempBlock.GetHash()
		BroadcastBlockHash(v, blockHash)
		v.MemPool.Clear()
		return 0
	}

	counter++
	return counter
}
