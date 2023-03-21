package types


type MemPool struct {
	transactions []*Transaction
}

func (m *MemPool) AddTransaction(tx *Transaction) {
	m.transactions = append(m.transactions, tx)
}

func (m *MemPool) Size() int {
	return len(m.transactions)
}

func NewMemPool() *MemPool {
	return &MemPool{}
}