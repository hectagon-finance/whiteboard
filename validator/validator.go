package validator

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

var Chan_1 = make(chan Msg, 100)

// amount of transactions to remove from mempool
var Chan_2 = make(chan int)

var DraftBlock = Block{} // Propose future block

var Port = ""

var Peers = []string{}

var Chain = Blockchain{}

var ReceivedBlockHash = make(map[string]string)

var MemPoolValidator = MemPool{}

type Msg struct {
	memPool MemPool
	heigt   int
}

type Validator struct {
	PublicKey    string
	PrivateKey   string
}

func (v *Validator) GetPublicKey() string {
	return v.PublicKey
}

func (v *Validator) GetPrivateKey() string {
	return v.PrivateKey
}

type Consensus struct {
	ReceivedMessage []map[string]interface{}
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

		HandleMessage(v, message)
		k := <- Chan_2
		MemPoolValidator.Transactions = MemPoolValidator.Transactions[k:]
	}
}
