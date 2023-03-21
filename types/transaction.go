package types

import (
	"github.com/hectagon-finance/whiteboard/crypto"
	"github.com/hectagon-finance/whiteboard/utils"
)

type Transaction interface {
	// Get the transaction's id
	Id() string

	// Get the transaction's public key
	PublicKey() crypto.PublicKey
	// Get the transaction's timestamp

	// Get the transaction's signature
	Signature() crypto.Signature

	// Get the transaction's hash
	Data() []byte
}

type transaction struct {
	transactionId string
	publicKey     crypto.PublicKey
	timestamp     int64
	signature     crypto.Signature
	data          []byte
}

func (t *transaction) Id() string {
	return t.transactionId
}

func (t *transaction) Timestamp() int64 {
	return t.timestamp
}

func (t *transaction) Signature() crypto.Signature {
	return t.signature
}

func (t *transaction) Data() []byte {
	return t.data
}

func (t *transaction) PublicKey() crypto.PublicKey {
	return t.publicKey
}

// func NewTransaction(transactionId string, publicKey string, signature string, hash string) Transaction {
// 	return &transaction{
// 		transactionId: transactionId,
// 		publicKey:     publicKey,
// 		timestamp:     time.Now().UnixNano(),
// 		signature:     signature,
// 		hash:          hash,
// 	}
// }

func NewTransaction(publicKey crypto.PublicKey, signature crypto.Signature, data []byte) Transaction {
	id := utils.RandString(9)
	return &transaction{
		transactionId: id,
		publicKey:     publicKey,
		signature:     signature,
		data:          data,
	}
}

