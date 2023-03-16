package types

import (
	"time"
)

type  Block  interface  {
	// Get the block's id
	Id ()  string

	// Get the block's timestamp
	Timestamp ()  int64

	// Get the block's hash
	Hash ()  string

	// Get the block's previous hash
	PreviousHash ()  string

	// Get the block's transactions
	Transactions () [] Transaction

	// Get the block's validator
	Validator ()  Validator

	// Get the block's signature
	Signature ()  string

}

type  block  struct  {
	blockId  string
	timestamp  int64
	hash  string
	previousHash  string
	transactions  [] Transaction
	validator  Validator
	signature  string
}

func  (b * block )  Id ()  string  {
	 return  b.blockId
}

func  (b * block )  Timestamp ()  int64  {
	 return  b.timestamp
}

func  (b * block )  Hash ()  string  {
	 return  b.hash
}

func  (b * block )  PreviousHash ()  string  {
	 return  b.previousHash
}

func  (b * block )  Transactions () [] Transaction  {
	 return  b.transactions
}

func  (b * block )  Validator ()  Validator  {
	 return  b.validator
}

func  (b * block )  Signature ()  string  {
	 return  b.signature
}

func  NewBlock (blockId string, hash string, previousHash string, transactions [] Transaction, validator Validator, signature string)  Block  {
	return & block  {
		blockId: blockId,
		timestamp: time.Now().UnixNano(),
		hash: hash,
		previousHash: previousHash,
		transactions: transactions,
		validator: validator,
		signature: signature,
	}
}

func NewGenesisBlock (validator Validator)  Block  {
	return & block  {
		blockId: "0",
		timestamp: time.Now().UnixNano(),
		hash: "0",
		previousHash: "0",
		transactions: [] Transaction {},
		validator: validator,
		signature: "0",
	}
}

func NewRandomBlock (validator Validator)  Block  {
	return & block  {
		blockId: "0",
		timestamp: time.Now().UnixNano(),
		hash: "0",
		previousHash: "0",
		transactions: [] Transaction {},
		validator: validator,
		signature: "0",
	}
}

