package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
	. "github.com/hectagon-finance/whiteboard/validator/start"
)

func main() {
	if len(os.Args) >= 2 {
		current_validator_id := os.Args[1]

		a, err := strconv.Atoi(current_validator_id)
		if err != nil {
			// ... handle error
			panic(err)
		}

		// Create two validators
		v := NewValidator(a)

		// Start the validators
		Start(&v)

		// Wait for a few seconds to let the validators establish connections
		time.Sleep(100 * time.Second)

		// Stop the validators
		Stop(&v)
	}
}

func NewValidator(port int) Validator {
	var blockChain Blockchain
	blockChain = NewBlockchain()

	validatorId := strconv.Itoa(port)
	publicKey := "public-key"
	privateKey := "private-key"
	memPool := NewMemPool()

	peers := []string{"8080"}

	if port != 8080 {
		peers = append(peers, strconv.Itoa(port))
	}

	return Validator{
		ValidatorId: validatorId,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		Blockchain:  blockChain,
		MemPool:     memPool,
		Balance:     0,
		Stake:       0,
		Status:      "inactive",
		Port:        port,
		Clients:     make(map[*websocket.Conn]bool),
		Peers:       peers,
	}
}
