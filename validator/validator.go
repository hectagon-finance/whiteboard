package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils"
)

var Chan_1 = make(chan Msg, 100)

// amount of transactions to remove from mempool
var Chan_2 = make(chan int)

var DraftBlock = types.Block{} // Propose future block

var Port = ""

var Peers = []string{}

var Chain = types.Blockchain{}

var ReceivedBlockHash = make(map[string]string)

var MemPoolValidator = types.MemPool{}

var ShouldReceiveTxFromPeer = true

// type Chan1Message struct {
// 	Time  bool
// 	Msg   Msg
// }

type Msg struct {
	memPool types.MemPool
	heigt   int
}

type Validator struct {
	PublicKey  string
	PrivateKey string
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
		// log.Printf("recv: %s", message)

		HandleMessage(v, message)
		k := <-Chan_2
		MemPoolValidator.Transactions = MemPoolValidator.Transactions[k:]
	}
}

func HandleMessage(v *Validator, msg []byte) {
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error unmarshaling the message:", err)
		return
	}

	if message["from"].(string) != "client" {
		exist := false
		for _, peer := range Peers {
			if peer == message["from"].(string) {
				exist = true
			}
		}
		if !exist {
			Peers = append(Peers, message["from"].(string))
			log.Println("Add peer", message["from"].(string))
		}
	}

	switch message["type"].(string) {

	case "sync all request":
		SyncAllRequestHandler(message)
	case "sync all response":
		SyncAllResponseHandler(message)
	case "sync draft block request":
		SyncBlockDraftRequestHandler(message)
	case "sync draft block response":
		SyncBlockDraftResponseHandler(message)
	case "transaction":
		HandleTypeTransaction(v, message)
	case "blockHash":
		ReceivedBlockHash[message["from"].(string)] = message["blockHash"].(string)
		fmt.Println("handler.go ** before checkCondition")
		produceBlock()
	default:
		fmt.Println("Default")
	}
}

func produceBlock() {
	log.Println("Check condition")
	log.Println("DraftBlock.Hash", DraftBlock.Hash)
	if DraftBlock.Hash == [32]byte{} {
		log.Println("DraftBlock.Hash != [32]byte{}")
		return
	}
	totalReceived := 0.0
	BlockHashCounter := make(map[string]int)
	max := 0
	winner := ""
	finalHash := ""

	blockHashStr := utils.Byte32toStr(DraftBlock.Hash)

	ReceivedBlockHash[Port] = blockHashStr
	for peer := range ReceivedBlockHash {
		BlockHashCounter[ReceivedBlockHash[peer]]++
		if BlockHashCounter[ReceivedBlockHash[peer]] > max {
			max = BlockHashCounter[ReceivedBlockHash[peer]]
			winner = peer
			finalHash = ReceivedBlockHash[peer]
		}
		totalReceived++
	}
	log.Println("Total received", totalReceived)
	log.Println("Max", max)
	log.Println("Percent", float64(len(Peers))*0.7)
	log.Println("Final hash", finalHash, "blockHashStr", blockHashStr)
	if totalReceived >= float64(len(Peers))*0.7 {
		// add to chain
		if blockHashStr == finalHash {
			preBlockHash := Chain.LastBlock().Hash
			Chain.CreateBlock(DraftBlock.Height, preBlockHash, DraftBlock.Transactions)
			Chan_Block <- DraftBlock
			ShouldReceiveTxFromPeer = true
			Chain.Print()
		} else {
			// sync from winner
			log.Println("Sync Draft Block from winner", winner)

			u := url.URL{Scheme: "ws", Host: "localhost:" + winner, Path: "/ws"}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			defer conn.Close()
			if err != nil {
				log.Fatal("dial:", err)
			}
			msg := map[string]interface{}{
				"type":   "sync draft block request",
				"from":   Port,
				"to":     winner,
				"height": DraftBlock.Height,
			}
			message, err := json.Marshal(msg)
			if err != nil {
				log.Fatal("marshal:", err)
			}
			conn.WriteMessage(websocket.TextMessage, message)
		}

		// reset
		ReceivedBlockHash = make(map[string]string)
		DraftBlock = types.Block{}
	}
}
