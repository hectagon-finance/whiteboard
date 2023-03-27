package validator

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

type Validator struct {
	Addr 		 string
	PublicKey    string
	PrivateKey   string
	Blockchain   Blockchain
	MemPool      MemPool
	LastBlock    Block
	Port         string
	Peers        []string
}

func (v *Validator) GetPublicKey() string {
	return v.PublicKey
}

func (v *Validator) GetPrivateKey() string {
	return v.PrivateKey
}

func (v *Validator) GetMemPool() MemPool {
	return v.MemPool
}

func (v *Validator) GetLastBlock() Block {
	return v.LastBlock
}

func (v *Validator) GetPort() string {
	return v.Port
}

type Consensus struct {
	receivedMessage []map[string]interface{}
}

func (v *Validator) Serve(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	} 
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Handle the message coming
		HandleMessage(v, message)

	}
}

// func AddMessage(v *Validator, message map[string]interface{}) {

// 	b := v.Consensus
// 	b.receivedMessage = append(b.receivedMessage, message)

// 	totalMessage := 0
// 	blockHashCounter := make(map[string]int)
// 	for _, blockHash := range b.receivedMessage {
// 		blockHashCounter[blockHash["blockHash"].(string)]++
// 		totalMessage++
// 	}
// 	handleConsensus(v, blockHashCounter, totalMessage)

// }

// func handleConsensus(v *Validator, blockHashCounter map[string]int, totalMessage int) {
// 	for _, count := range blockHashCounter {
// 		if float64(count)/float64(totalMessage) >= 0.6 {
// 			preBlockHash := v.Blockchain.LastBlock().Hash
// 			v.Blockchain.CreateBlock(v.TempBlock.Height, preBlockHash, v.TempBlock.Transactions)
// 			v.Consensus = Consensus{}
// 			fmt.Printf("This is blockchain of %s \n", v.ValidatorId)
// 			v.Blockchain.Print()
// 		}
// 	}
// }
