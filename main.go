package main

import (
	"log"
	"os"

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
		go ClientHandler(&v, is_genesis)

		go BroadcastBlockHash()

		// Wait for a few seconds to let the validators establish connections
		for {
			
		}
	}
}
// ./main 8080 genesis ; ./main 9000 8080
func NewValidator(port string, is_genenis string) Validator {
	Port = port
	log.Println("Port: ", Port)
	Peers = append(Peers, port)

	if is_genenis == "genesis" {
		// genesis validator	
		Chain = NewBlockchain()
	}else {
		// sync with other is_genenis (port)
		Peers = append(Peers, is_genenis)
	}
	
	publicKey := "public-key"
	privateKey := "private-key"
	MemPoolValidator = NewMemPool()

	return Validator{
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
	}
}
