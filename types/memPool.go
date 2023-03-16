package types

type MemPool interface {
	// Add a transaction to the mempool
	AddTransaction(tx Transaction)

	// Remove a transaction from the mempool
	RemoveTransaction(tx Transaction)

	// Get the transactions in the mempool
	GetTransactions() []Transaction

	// Get the size of the mempool
	Size() int
}

type memPool struct {
	transactions []Transaction
}

func (m *memPool) AddTransaction(tx Transaction) {
	m.transactions = append(m.transactions, tx)
}

func (m *memPool) RemoveTransaction(tx Transaction) {
	for i, t := range m.transactions {
		if t.Id() == tx.Id() {
			m.transactions = append(m.transactions[:i], m.transactions[i+1:]...)
		}
	}
}

func (m *memPool) GetTransactions() []Transaction {
	return m.transactions
}

func (m *memPool) Size() int {
	return len(m.transactions)
}

func NewMemPool() MemPool {
	return &memPool{}
}