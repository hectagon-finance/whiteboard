package types

import (
	"time"
)

type Transaction interface {
	// Get the transaction's id
	Id() string

	// Get the transaction's public key
	PublicKey() string

	// Get the transaction's timestamp
	Timestamp() int64

	// Get the transaction's signature
	Signature() string

	// Get the transaction's hash
	Hash() string
}

type transaction struct {
	transactionId	string
	publicKey 		string
	timestamp 		int64
	signature 		string
	hash      		string
}

func (t *transaction) Id() string {
	return t.transactionId
}


func (t *transaction) Timestamp() int64 {
	return t.timestamp
}

func (t *transaction) Signature() string {
	return t.signature
}

func (t *transaction) Hash() string {
	return t.hash
}

func (t *transaction) PublicKey() string {
	return t.publicKey
}

func NewTransaction(transactionId string, publicKey string, signature string, hash string) Transaction {
	return &transaction{
		transactionId: transactionId,
		publicKey: publicKey,
		timestamp: time.Now().UnixNano(),
		signature: signature,
		hash: hash,
	}
}