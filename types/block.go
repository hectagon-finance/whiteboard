package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	height       int
	hash         [32]byte
	previousHash [32]byte
	transactions []Transaction
}

func (b *Block) Id() int {
	return b.height
}

func (b *Block) Hash() [32]byte {
	return b.hash
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Transactions() []Transaction {
	return b.transactions
}

func NewBlock(height int, previousHash [32]byte, transactions []Transaction) *Block {
	b := new(Block)
	b.height = height
	b.previousHash = previousHash
	b.transactions = transactions
	b.hash = b.calculateHash()
	return b
}

func (b *Block) calculateHash() [32]byte {
	m, _ := json.Marshal(struct {
		Height       int            `json:"height"`
		PreviousHash [32]byte       `json:"previousHash"`
		Transactions []Transaction `json:"transactions"`
	}{
		Height:       b.height,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
	fmt.Println(string(m))
	// hash height, previousHash, transactions
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {

	fmt.Printf("previous_hash   %x\n", b.previousHash)
	fmt.Printf("hash   %x\n", b.hash)
	for _, tx := range b.transactions {
		fmt.Println("Transaction:", tx.Data)	
	}
}
