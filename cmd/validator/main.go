package main

import (
	"os"
	"time"

	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func main() {
	if len(os.Args) >= 2 {
		current_validator_id := os.Args[1]
		is_genesis := os.Args[2]

		// Create two validators
		v := NewValidator(current_validator_id, is_genesis)

		// Start the validators
		ClientHandler(&v, is_genesis)

		// Wait for a few seconds to let the validators establish connections
		time.Sleep(100 * time.Second)
	}
}
// ./main 8080 genesis ; ./main 9000 8080
func NewValidator(port string, is_genenis string) Validator {
	var blockChain Blockchain
	peers := []string{}
	peers = append(peers, port)
	if is_genenis == "genesis" {
		// genesis validator	
		blockChain = NewBlockchain()
	}else {
		// sync with other is_genenis (port)
		peers = append(peers, is_genenis)
	}
	
	publicKey := "public-key"
	privateKey := "private-key"
	memPool := NewMemPool()

	return Validator{
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		Blockchain:  blockChain,
		MemPool:     memPool,
		Port:        port,
		Peers:       peers,
	}
}
