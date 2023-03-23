package types

import "encoding/json"

type MemPool struct {
	Transactions []Transaction
}

type MemPoolEncode struct {
	Transactions []TransactionEncode
}

func (m *MemPool) AddTransaction(tx Transaction) {
	m.Transactions = append(m.Transactions, tx)
}

func (m *MemPool) GetTransactions() []Transaction {
	return m.Transactions
}

func (m *MemPool) Size() int {
	return len(m.Transactions)
}

func (m *MemPool) Clear() {
	m.Transactions = []Transaction{}
}

func NewMemPool() MemPool {
	return MemPool{}
}

func (m *MemPool) Encode() ([]byte, error) {
	memPoolencode := MemPoolEncode{}
	for _, tx := range m.Transactions {
		memPoolencode.Transactions = append(memPoolencode.Transactions, tx.Tranform())
	}

	memByte, err := json.Marshal(memPoolencode)
	if err != nil {
		return nil, err
	}
	return memByte, nil
}

func DecodeMempool(memBye []byte) (MemPool, error) {
	memPoolencode := MemPoolEncode{}
	err := json.Unmarshal(memBye, &memPoolencode)
	if err != nil {
		return MemPool{}, err
	}

	memPool := MemPool{}
	for _, tx := range memPoolencode.Transactions {
		memPool.Transactions = append(memPool.Transactions, tx.Tranform())
	}

	return memPool, nil
}
