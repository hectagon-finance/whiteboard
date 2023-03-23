package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Blockchain struct {
	Chain []Block
}

type BlockchainEncode struct {
	Chain []BlockEncode
}

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

func (bc *Blockchain) Encode() ([]byte, error) {
	blockchainEncode := BlockchainEncode{}
	for _, block := range bc.Chain {
		blockEncode, err := block.Tranform()
		if err != nil {
			return nil, err
		}
		blockchainEncode.Chain = append(blockchainEncode.Chain, blockEncode)
	}
	
	blockchainByte, err := json.Marshal(blockchainEncode)
	if err != nil {
		return nil, err
	}

	return blockchainByte, nil
}

func DecodeBlockchain(blockchainByte []byte) (Blockchain, error) {
	blockchainEncode := BlockchainEncode{}
	err := json.Unmarshal(blockchainByte, &blockchainEncode)
	if err != nil {
		return Blockchain{}, err
	}

	blockchain := Blockchain{}
	for _, blockEncode := range blockchainEncode.Chain {
		block, err := BlockEncodeTransform(blockEncode)
		if err != nil {
			return Blockchain{}, err
		}
		blockchain.Chain = append(blockchain.Chain, block)
	}

	return blockchain, nil
}