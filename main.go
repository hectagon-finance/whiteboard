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

		a, err := strconv.Atoi(current_validator_id)
		if err != nil {
			// ... handle error
			panic(err)
		}

		// Create two validators
		v := types.NewValidator(a)

		// Start the validators
		v.Start()

		// Wait for a few seconds to let the validators establish connections
		time.Sleep(100 * time.Second)

		// Stop the validators
		v.Stop()
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