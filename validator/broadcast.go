package validator

import (

	. "github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils"
)

func BroadcastPeer() {

	message := map[string]interface{}{
		"type":        "peer",
		"from":        Port,
		"message":     Peers,
	}

	ConnectAndSendMessage(message)
}

func BroadcastTransaction(tx Transaction) {
	publicKey := tx.PublicKey
	publicKeyStr := publicKey.PublicKeyStr()

	signature := tx.Signature
	signatureStr := signature.SignatureStr()
 
	blockHashStr := utils.Byte32toStr(Chain.LastBlock().Hash)



	message := map[string]interface{}{
		"type":          "transaction",
		"from":          Port,
		"transactionId": tx.Id(),
		"publicKey":     publicKeyStr,
		"signature":     signatureStr,
		"data":          string(tx.Data),
		"latestBlockHash":   blockHashStr,
	}

	ConnectAndSendMessage(message)
}
