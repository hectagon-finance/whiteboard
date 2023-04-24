package validator

import (
	"fmt"
	"log"

	"github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils"
	"github.com/hectagon-finance/whiteboard/utils/crypto"
)

func HandleTypeTransaction(v *Validator, message map[string]interface{}) {
	var lastBlockHashFromMessage string
	publicKeyStr := message["publicKey"].(string)
	signatureStr := message["signature"].(string)

	publicKey := crypto.PublicKeyFromString(publicKeyStr)
	signature := crypto.SignatureFromString(signatureStr)

	data := message["data"].(string)

	tx := types.Transaction{
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
			if !ShouldReceiveTxFromPeer {
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
			// Chan_1 <- Chan1Message {
			// 	Msg : Msg{MemPoolValidator, Chain.LastBlock().Height+1},
			// 	Time: false,
			// }

			Chan_1 <- Msg{MemPoolValidator, Chain.LastBlock().Height + 1}

		} else {
			fmt.Printf("Validator %s: Already have that transaction\n", Port)
		}
	} else {
		fmt.Print("Invalid")
		fmt.Printf("Validator %s: Invalid transaction received from %s: %s\n", Port, message["from"].(string), tx.Id())
	}
}

func checkTransaction(v *Validator, tx types.Transaction, lastBlockHash string) bool {
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

func BroadcastTransaction(tx types.Transaction) {
	publicKey := tx.PublicKey
	publicKeyStr := publicKey.PublicKeyStr()

	signature := tx.Signature
	signatureStr := signature.SignatureStr()

	blockHashStr := utils.Byte32toStr(Chain.LastBlock().Hash)

	message := map[string]interface{}{
		"type":            "transaction",
		"from":            Port,
		"transactionId":   tx.Id(),
		"publicKey":       publicKeyStr,
		"signature":       signatureStr,
		"data":            string(tx.Data),
		"latestBlockHash": blockHashStr,
	}

	ConnectAndSendMessage(message)
}
