package validator

import (
	"encoding/json"
	"fmt"

	. "github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils/crypto"
)

func HandleMessage(v *Validator, msg []byte) {
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error unmarshaling the message:", err)
		return
	}

	switch message["type"].(string) {
	
	case "sync":
		fmt.Println("Sync message received")
		peer_pool_size := message["memPoolSize"].(float64)
		if v.MemPool.Size() > 0 {
			
			if peer_pool_size == 0 {
				fmt.Println("Need to sync mempool for validator", message["Port"].(string))
				validator := v
				validator.Peers = []string{validator.Port, message["Port"].(string)}
				Sync(validator)
			}
		} else {
			if peer_pool_size > 0 {
				fmt.Println(v.MemPool.Size(), peer_pool_size)
				memPool, err := DecodeMempool([]byte(message["memPool"].(string)))
				if err != nil {
					panic(err)
				}
				v.MemPool = memPool

				blockchain, err := DecodeBlockchain([]byte(message["blockchain"].(string)))
				if err != nil {
					panic(err)
				}

				v.Blockchain = blockchain

				fmt.Println("Sync mempool and blockchain")
				
				fmt.Printf("Validator %s: Mempool size: %d\n", v.Port, v.MemPool.Size())
				
				v.Blockchain.Print()

			}
		}

	case "transaction":

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
			if checkTransaction(v, tx) {
				fmt.Printf("Validator %s: Valid transaction received from %s %s: %s\n", v.Port, message["from"].(string), message["Port"].(string), tx.Id())
				BroadcastTransaction(v, tx)
				v.MemPool.AddTransaction(tx)
				fmt.Printf("Validator %s: Adding new transaction to mempool: %s\n", v.Port, tx.Id())
				fmt.Printf("Validator %s: Mempool size: %d\n", v.Port, v.MemPool.Size())
			} else {
				fmt.Printf("Validator %s: Already have that transaction\n", v.Port)
			}
		} else {
			fmt.Print("Invalid")
			fmt.Printf("Validator %s: Invalid transaction received from %s %s: %s\n", v.Port, message["from"].(string), message["Port"].(string), tx.Id())
		}
	case "peer":
		addPeer(v, message["message"].([]interface{}))

	case "blockHash":
			// AddMessage(v, message)

	default:
		fmt.Println("Default")
	}

}

func addPeer(v *Validator, peers []interface{}) {
	for _, peer := range peers {
		alreadyhave := false
		for _, p := range v.Peers {
			if p == peer.(string) {
				alreadyhave = true
				break
			}
		}

		if !alreadyhave {
			fmt.Println("Validator", v.Port, ": Adding new peer", peer.(string))
			v.Peers = append(v.Peers, peer.(string))
		}

	}
}

func checkTransaction(v *Validator, tx Transaction) bool {
	for i := range v.MemPool.GetTransactions() {
		if v.MemPool.GetTransactions()[i].Id() == tx.Id() {
			return false
		}
	}
	return true
}
