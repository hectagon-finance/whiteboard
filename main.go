package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/onspeedhp/project/types"
)

func main() {
	if len(os.Args) >= 2{
		current_validator_id := os.Args[1]
		peer_validator_id := os.Args[2]

		a, err := strconv.Atoi(current_validator_id)
		b, err := strconv.Atoi(peer_validator_id)
		if err != nil {
			// ... handle error
			panic(err)
		}

		// Create two validators
		v1 := types.NewValidator(a)
		v2 := types.NewValidator(b)

		// Add each other as peers
		v1.AddPeer(peer_validator_id)
		v2.AddPeer(current_validator_id)

		// Start the validators
		v1.Start()

		// Wait for a few seconds to let the validators establish connections
		time.Sleep(100 * time.Second)

		// Stop the validators
		v1.Stop()
	} else {
		// Create a fake transaction
		tx := types.FakeTransaction()

		// Send the fake transaction to validator 1
		sendFakeTransaction("8080", tx)
	}
}

func sendFakeTransaction(validatorId string, tx types.Transaction) {
	u := url.URL{Scheme: "ws", Host: "localhost:" + validatorId, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to the validator:", err)
		return
	}
	defer conn.Close()

	message := map[string]interface{}{
		"type":          "transaction",
		"from": 		 "client",
		"validatorId":   "fake-client",
		"transactionId": tx.Id(),
		"publicKey":     tx.PublicKey(),
		"timestamp":     tx.Timestamp(),
		"signature":     tx.Signature(),
		"hash":          tx.Hash(),
	}

	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling the message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}