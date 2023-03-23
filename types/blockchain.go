package types

import (
	"fmt"
	"strings"
)

type Blockchain struct {
	// memPool MemPool
	Chain []Block
}

// func (bc *Blockchain) AddTransaction(tx *Transaction) {
// 	bc.memPool.AddTransaction(tx)
// }

func NewBlockchain() Blockchain {
	b := &Block{}
	bc := *new(Blockchain)
	transactions := []Transaction{}
	bc.CreateBlock(0, b.Hash, transactions)
	return bc
}

func (bc *Blockchain) CreateBlock(height int, previousHash [32]byte, transactions []Transaction) {
	b := NewBlock(height, previousHash, transactions)
	bc.Chain = append(bc.Chain, b)

	// return b
}

func (bc *Blockchain) LastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}
