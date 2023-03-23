package start

import (
	"encoding/hex"

	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func BroadcastMempool(v *Validator) {

	// encode validator mempool
	memByte, err := v.MemPool.Encode()
	if err != nil {
		panic(err)
	}

	message := map[string]interface{}{
		"type":        "memPool",
		"validatorId": v.Id(),
		"memPoolSize": v.MemPool.Size(),
		"memPool":     string(memByte),
		"message":     "Hello, I'm validator " + v.Id(),
	}

	ConnectAndSendMessage(v, message)
}

func BroadcastPeer(v *Validator) {

	message := map[string]interface{}{
		"type":        "peer",
		"from":        v.ValidatorId,
		"validatorId": v.ValidatorId,
		"message":     v.Peers,
	}

	ConnectAndSendMessage(v, message)
}

func BroadcastTransaction(v *Validator, tx Transaction) {
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

func BroadcastBlockHash(v *Validator, blockHash [32]byte) {
	blockHashSlice := blockHash[:]
	blockHashStr := hex.EncodeToString(blockHashSlice)

	message := map[string]interface{}{
		"type":        "blockHash",
		"validatorId": v.ValidatorId,
		"blockHash":   blockHashStr,
	}

	ConnectAndSendMessage(v, message)
}
