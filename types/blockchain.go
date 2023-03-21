package types

import (
	"fmt"
	"strings"
)

type Blockchain struct {
	// memPool MemPool
	chain []*Block
}

// func (bc *Blockchain) AddTransaction(tx *Transaction) {
// 	bc.memPool.AddTransaction(tx)
// }

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	transactions := []*Transaction{}
	bc.CreateBlock(0, b.hash, transactions)
	return bc
}

func (bc *Blockchain) CreateBlock(height int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := NewBlock(height, previousHash, transactions)
	fmt.Println("createblock", b)
	bc.chain = append(bc.chain, b)

	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}
