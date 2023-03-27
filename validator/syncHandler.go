package validator

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

func SyncBlockDraftRequestHandler(message map[string]interface{}){
	height := message["height"].(float64)
	result := Blockchain{}
	
	if DraftBlock.Height == int(height) {
		result.Chain = append(result.Chain, DraftBlock)
	} else {
		result = Chain
	}

	u := url.URL{Scheme: "ws", Host: "localhost:" + message["from"].(string), Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil) 
	defer conn.Close()
	if err != nil {
		log.Fatal("dial:", err)
	}

	resultBytes, err := result.Encode()
	if err != nil {
		log.Fatal("encode:", err)
	}

	msg := map[string]interface{}{
		"type": "sync draft block response",
		"from": Port,
		"result": resultBytes,
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("marshal:", err)
	}
	conn.WriteMessage(websocket.TextMessage, msgByte)
}

func SyncBlockDraftResponseHandler(message map[string]interface{}){
	chainbyte := message["result"].(string)
	blockchain, err := DecodeBlockchain([]byte(chainbyte))

	if err != nil {
		log.Fatal("decode:", err)
	}

	if len(blockchain.Chain) == 1 {
		newBlock := blockchain.Chain[0]
		Chain.CreateBlock(newBlock.Height, newBlock.PreviousHash, newBlock.Transactions)
	} else {
		Chain = blockchain
	}
}

func SyncAllRequestHandler(message map[string]interface{}){
	Chain.Print()
	// TODO: add new mempool to everyone in the network
	existed := false
	for _, peer := range Peers {
		if peer == message["from"].(string) {
			existed = true
		}
	}
	if !existed {
		Peers = append(Peers, message["from"].(string))
	}

	chainBytes, err := Chain.Encode()
	if err != nil {
		log.Fatal("encode:", err)
	}

	memByte, err := MemPoolValidator.Encode()
	if err != nil {
		panic(err)
	}

	u := url.URL{Scheme: "ws", Host: "localhost:" + message["from"].(string), Path: "/ws"}
	log.Print("Send response to: ", message["from"].(string))
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil) 
	defer conn.Close()
	if err != nil {
		log.Fatal("dial:", err)
	}

	msg := map[string]interface{}{
		"type": "sync all response",
		"from": Port,
		"to"  : message["from"].(string),
		"peers": Peers,
		"memPool": string(memByte),
		"chain": string(chainBytes),
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("marshal:", err)
	}


	conn.WriteMessage(websocket.TextMessage, msgByte)

}

func SyncAllResponseHandler(message map[string]interface{}){
	peers := message["peers"].([]interface{})
	AddPeer(peers)
	memPoolStr := message["memPool"].(string)
	chainStr := message["chain"].(string)

	// update peers
	// Peers = peers

	// update memPool
	memPool, err := DecodeMemPool([]byte(memPoolStr))
	if err != nil {
		panic(err)
	}

	MemPoolValidator.Transactions = memPool.Transactions

	// update chain
	chain, err := DecodeBlockchain([]byte(chainStr))
	if err != nil {
		panic(err)
	}
	Chain = chain
	Chain.Print()
}