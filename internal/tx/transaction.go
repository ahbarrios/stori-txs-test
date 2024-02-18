// Transaction is the main package of this software.
// It will contain every important aspect to handle bussiness workflows related
// with finantial Transactions within the system.
package tx

import (
	"errors"
	"time"
)

// NoRecordError it will describe that you consume all records from the source
// is up to the source to handle this error
var NoRecordError = errors.New("no more records available in source")

// Reader any source will implement this to extract transactions from it one at a time.
// The style of the Reader interface is iterator-like or streaming. So the Transaction processor
// will keep reading using the Read() method from the source until the [internal/tx/NoRecordError] gets returned.
type Reader interface {
	Read() (Transaction, error)
}

// Transaction is the main Domain model for this product.
// It will contain all the important information of a money transfer transaction for the given system.
type Transaction struct {
	ID     int
	Date   time.Time
	Amount float64
}
