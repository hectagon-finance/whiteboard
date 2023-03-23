package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
)

type Block struct {
	Height       int
	Hash         [32]byte
	PreviousHash [32]byte
	Transactions []Transaction
}

type BlockEncode struct {
	Height       string
	Hash         string
	PreviousHash string
	Transactions []TransactionEncode
}

func (b *Block) Id() int {
	return b.Height
}

func (b *Block) GetHash() [32]byte {
	return b.Hash
}

func (b *Block) GetPreviousHash() [32]byte {
	return b.PreviousHash
}

func (b *Block) GetTransactions() []Transaction {
	return b.Transactions
}

func NewBlock(Height int, PreviousHash [32]byte, Transactions []Transaction) Block {
	b := *new(Block)
	b.Height = Height
	b.PreviousHash = PreviousHash
	b.Transactions = Transactions
	b.Hash = b.calculateHash()
	return b
}

func (b *Block) calculateHash() [32]byte {
	m, _ := json.Marshal(struct {
		Height       int           `json:"height"`
		PreviousHash [32]byte      `json:"previousHash"`
		Transactions []Transaction `json:"transactions"`
	}{
		Height:       b.Height,
		PreviousHash: b.PreviousHash,
		Transactions: b.Transactions,
	})
	// hash height, previousHash, transactions
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {

	fmt.Printf("previous_hash   %x\n", b.PreviousHash)
	fmt.Printf("hash   %x\n", b.Hash)
	for _, tx := range b.Transactions {
		fmt.Println("Transaction:", string(tx.Data))
	}
}

func (b *Block) Tranform() (BlockEncode, error) {
	blockEncode := BlockEncode{}
	for _, tx := range b.Transactions {
		blockEncode.Transactions = append(blockEncode.Transactions, tx.Tranform())
	}

	return blockEncode, nil
}

func BlockEncodeTransform(blockEncode BlockEncode) (Block, error) {

	block := Block{}

	block.Height, _ = strconv.Atoi(blockEncode.Height)
	copy(block.Hash[:], blockEncode.Hash)
	copy(block.PreviousHash[:], blockEncode.PreviousHash)
	for _, txEncode := range blockEncode.Transactions {
		block.Transactions = append(block.Transactions, txEncode.Tranform())
	}

	return block, nil
}