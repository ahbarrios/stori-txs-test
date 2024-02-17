// Transaction is the main package of this software.
// It will contain every important aspect to handle bussiness workflows related
// with finantial Transactions within the system.
package tx

import "time"

// Source any source will implement this to extract transactions from it one at a time.
// The style of the Source interface is iterator-like or streaming. So the Transaction processor
// will keep reading using the Read() method from the source until the Done() is true
type Source interface {
	Read() (Transaction, error)
	Done() bool
}

// Transaction is the main Domain model for this product.
// It will contain all the important information of a money transfer transaction for the given system.
type Transaction struct {
	ID     int
	Date   time.Time
	Amount float64
}
