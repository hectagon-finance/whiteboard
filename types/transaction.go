package types

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hectagon-finance/whiteboard/utils/crypto"
)

type Transaction struct {
	TransactionId string
	PublicKey     crypto.PublicKey
	Signature     crypto.Signature
	Data          []byte
}

type TransactionEncode struct {
	TransactionId string
	PublicKey     string
	Signature     string
	Data          string
}

func (t *Transaction) Id() string {
	return t.TransactionId
}

func (t *Transaction) GetPublicKey() crypto.PublicKey {
	return t.PublicKey
}

func (t *Transaction) GetSignature() crypto.Signature {
	return t.Signature
}

func (t *Transaction) GetData() []byte {
	return t.Data
}

func NewTransaction(publicKey crypto.PublicKey, signature crypto.Signature, data []byte) Transaction {
	id := strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(rand.Intn(1000000))
	return Transaction{
		TransactionId: id,
		PublicKey:     publicKey,
		Signature:     signature,
		Data:          data,
	}
}

func (t *Transaction) Tranform() TransactionEncode {
	return TransactionEncode{
		TransactionId: t.TransactionId,
		PublicKey:     t.PublicKey.PublicKeyStr(),
		Signature:     t.Signature.SignatureStr(),
		Data:          string(t.Data),
	}
}

func (t *TransactionEncode) Tranform() Transaction {
	return Transaction{
		TransactionId: t.TransactionId,
		PublicKey:     *crypto.PublicKeyFromString(t.PublicKey),
		Signature:     *crypto.SignatureFromString(t.Signature),
		Data:          []byte(t.Data),
	}
}
