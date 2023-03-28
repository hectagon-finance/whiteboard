package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils"
	"github.com/hectagon-finance/whiteboard/utils/crypto"
)

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
		var lastBlockHashFromMessage string
		publicKeyStr := message["publicKey"].(string)
		signatureStr := message["signature"].(string)

		publicKey := crypto.PublicKeyFromString(publicKeyStr)
		signature := crypto.SignatureFromString(signatureStr)

		data := message["data"].(string)

		tx := Transaction{
			TransactionId: message["transactionId"].(string),
			PublicKey:     *publicKey,
			Signature:     *signature,
			Data:          []byte(data),
		}

		if tx.Signature.Verify(*publicKey, tx.Data) {
			fmt.Printf("Validator %s: Valid transaction received from %s: %s\n", Port, message["from"].(string), tx.Id())
			
			if message["from"].(string) == "client" {
				lastBlockHashFromMessage = utils.Byte32toStr(Chain.LastBlock().Hash)
			} else {
				if ShouldReceiveTxFromPeer == false{
					return
				}
				lastBlockHashFromMessage = message["latestBlockHash"].(string)
			}

			if checkTransaction(v, tx, lastBlockHashFromMessage) {
				MemPoolValidator.AddTransaction(tx)
				BroadcastTransaction(tx)
				fmt.Printf("Validator %s: Adding new transaction to mempool: %s\n", Port, tx.Id())
				fmt.Printf("Validator %s: Mempool size: %d\n", Port, MemPoolValidator.Size())

				log.Println("Push to chan_1")
				Chan_1 <- Chan1Message {
					Msg : Msg{MemPoolValidator, Chain.LastBlock().Height+1},
					Time: false,
				}
				
			} else {
				fmt.Printf("Validator %s: Already have that transaction\n", Port)
			}
		} else {
			fmt.Print("Invalid")
			fmt.Printf("Validator %s: Invalid transaction received from %s: %s\n", Port, message["from"].(string), tx.Id())
		}
	case "blockHash":
		ReceivedBlockHash[message["from"].(string)] = message["blockHash"].(string)
		fmt.Println("handler.go ** before checkCondition")
		checkCondition()
	default:
		fmt.Println("Default")
	}
}

func checkCondition(){
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
	for peer, _ := range ReceivedBlockHash {
		BlockHashCounter[ReceivedBlockHash[peer]] ++
		if BlockHashCounter[ReceivedBlockHash[peer]] > max {
			max = BlockHashCounter[ReceivedBlockHash[peer]]
			winner = peer
			finalHash = ReceivedBlockHash[peer]
		}
		totalReceived ++
	}
	log.Println("Total received", totalReceived)
	log.Println("Max", max)
	log.Println("Percent", float64(len(Peers)) * 0.7)
	log.Println("Final hash", finalHash, "blockHashStr", blockHashStr)
	if totalReceived >= float64(len(Peers)) * 0.7 {
		// add to chain
		if blockHashStr == finalHash {
			preBlockHash := Chain.LastBlock().Hash
			Chain.CreateBlock(DraftBlock.Height, preBlockHash, DraftBlock.Transactions)
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
				"type": "sync draft block request",
				"from": Port,
				"to"  : winner,
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
		DraftBlock = Block{}
	}
}

func AddPeer(peers []interface{}) {
	for _, peer := range peers {
		alreadyhave := false
		for _, p := range Peers {
			if p == peer.(string) {
				alreadyhave = true
				break
			}
		}

		if !alreadyhave {
			fmt.Println("Validator", Port, ": Adding new peer", peer.(string))
			Peers = append(Peers, peer.(string))
		}

	}
}

func checkTransaction(v *Validator, tx Transaction, lastBlockHash string) bool {
	if lastBlockHash != utils.Byte32toStr(Chain.LastBlock().Hash) {
		return false
	}

	for i := range MemPoolValidator.GetTransactions() {
		if MemPoolValidator.GetTransactions()[i].Id() == tx.Id() {
			return false
		}
	}
	return true
}
