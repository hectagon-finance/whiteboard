package types


type MemPool struct {
	Transactions []Transaction
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

func NewMemPool() *MemPool {
	return &MemPool{}
}
