package start

import (
	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func CheckForAvailableValidators(v *ValidatorStruct) {

	message := map[string]interface{}{
		"type":         "hello",
		"validatorId":  v.Id(),
		"memPool":  	v.MemPool.Size(),
		"message":      "Hello, I'm validator " + v.Id(),
	}

	ConnectAndSendMessage(v, message)
}

func BroadcastPeer(v *ValidatorStruct) {

	message := map[string]interface{}{
		"type":        "peer",
		"from":        v.ValidatorId,
		"validatorId": v.ValidatorId,
		"message":     v.Peers,
	}

	ConnectAndSendMessage(v, message)
}

func BroadcastTransaction(v *ValidatorStruct, tx Transaction) {
	publicKey := tx.PublicKey
	publicKeyStr := publicKey.PublicKeyStr()

	signature := tx.Signature
	signatureStr := signature.SignatureStr()

	message := map[string]interface{}{
		"type":          "transaction",
		"from":          "client",
		"validatorId":   v.Id(),
		"transactionId": tx.Id(),
		"publicKey":     publicKeyStr,
		"signature":     signatureStr,
		"data":          string(tx.Data),
	}

	ConnectAndSendMessage(v, message)
}
