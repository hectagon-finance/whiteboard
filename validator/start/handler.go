package start

import (
	"encoding/json"
	"fmt"
	"github.com/hectagon-finance/whiteboard/utils/crypto"
	. "github.com/hectagon-finance/whiteboard/validator"
	. "github.com/hectagon-finance/whiteboard/types"
)

func HandleMessage(v *ValidatorStruct, msg []byte) {
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error unmarshaling the message:", err)
		return
	}

	switch message["type"].(string) {
		case "memPool":
			if v.MemPool.Size() == 0 {
				fmt.Println("Validator", v.ValidatorId, ": Received memPool from validator", message["validatorId"].(string))
				memByte := []byte(message["message"].(string))
				var memPool []Transaction
				err := json.Unmarshal(memByte, &memPool)
				if err != nil {
					fmt.Println("Error unmarshaling the message:", err)
					return
				}
				
			}
		case "hello":
			memPool_peer := int(message["memPool"].(float64))
			if memPool_peer == 0 {
				if v.MemPool.Size() == 0 {
				} else {
					memByte, err := json.Marshal(v.MemPool.Transactions)

					if err != nil {
						fmt.Println("Error marshaling the message:", err)
						return
					}
					message := map[string]interface{}{
						"type" : "memPool",
						"validatorId": v.ValidatorId,
						"message": string(memByte),
					}

					fmt.Println(memByte)
					fmt.Println("Vaidator", v.ValidatorId, ": Sending memPool to another validator")

					ConnectAndSendMessage(v, message)
				}
			}
			fmt.Println("Validator", v.ValidatorId, ": Received message from validator", message["validatorId"].(string), ":", message["message"].(string))
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
					fmt.Printf("Validator %s: Valid transaction received from %s %s: %s\n", v.ValidatorId, message["from"].(string), message["validatorId"].(string), tx.Id())
					BroadcastTransaction(v, tx)
					v.MemPool.AddTransaction(tx)
					fmt.Printf("Validator %s: Adding new transaction to mempool: %s\n", v.ValidatorId, tx.Id())
				} else {
					fmt.Printf("Validator %s: Already have that transaction\n", v.ValidatorId)
				}
			} else {
				fmt.Print("Invalid")
				fmt.Printf("Validator %s: Invalid transaction received from %s %s: %s\n", v.ValidatorId, message["from"].(string), message["validatorId"].(string), tx.Id())
			}
		case "peer":
			fmt.Printf("Validator %s: Valid peers array received from %s: %s\n", v.ValidatorId, message["validatorId"].(string), message["message"].([]interface{}))
			addPeer(v, message["message"].([]interface{}))
		default:
			fmt.Println("Default")
	}
}

func addPeer(v *ValidatorStruct, peers []interface{}) {
	for _, peer := range peers {
		alreadyhave := false
		for _, p := range v.Peers {
			if p == peer.(string) {
				alreadyhave = true
				break
			}
		}

		if !alreadyhave {
			fmt.Println("Validator", v.ValidatorId, ": Adding new peer", peer.(string))
			v.Peers = append(v.Peers, peer.(string))
		}

	}
}

func checkTransaction(v *ValidatorStruct, tx Transaction) bool {
	for i := range v.MemPool.GetTransactions() {
		if v.MemPool.GetTransactions()[i].Id() == tx.Id() {
			return false
		}

	}

	return true
}