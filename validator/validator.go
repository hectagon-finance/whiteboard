package validator

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

type Validator struct {
	ValidatorId  string
	PublicKey    string
	PrivateKey   string
	Blockchain   Blockchain
	MemPool      MemPool
	Consensus    Consensus
	TempBlock    Block
	Balance      int64
	Stake        int64
	Status       string
	LastBlock    Block
	Port         string
	HttpServer   *http.Server
	ClientsMutex sync.Mutex
	Clients      map[*websocket.Conn]bool
	Peers        []string
	StopServer   func()
}

func (v *Validator) Id() string {
	return v.ValidatorId
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

func (v *Validator) GetBalance() int64 {
	return v.Balance
}

func (v *Validator) GetStake() int64 {
	return v.Stake
}

func (v *Validator) GetStatus() string {
	return v.Status
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

func AddMessage(v *Validator, message map[string]interface{}) {

	b := v.Consensus
	b.receivedMessage = append(b.receivedMessage, message)

	totalMessage := 0
	blockHashCounter := make(map[string]int)
	for _, blockHash := range b.receivedMessage {
		blockHashCounter[blockHash["blockHash"].(string)]++
		totalMessage++
	}
	handleConsensus(v, blockHashCounter, totalMessage)

}

func handleConsensus(v *Validator, blockHashCounter map[string]int, totalMessage int) {
	for _, count := range blockHashCounter {
		if float64(count)/float64(totalMessage) >= 0.6 {
			preBlockHash := v.Blockchain.LastBlock().Hash
			v.Blockchain.CreateBlock(v.TempBlock.Height, preBlockHash, v.TempBlock.Transactions)
			v.Consensus = Consensus{}
			fmt.Printf("This is blockchain of %s \n", v.ValidatorId)
			v.Blockchain.Print()
		}
	}
}
