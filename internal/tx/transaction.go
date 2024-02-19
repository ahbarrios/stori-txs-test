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

// Producer any source will implement this to extract transactions from it one at a time.
// The style of the Producer interface is iterator-like or streaming. So the Transaction processor
// will keep reading using the Get() method from the source until the [internal/tx/NoRecordError] gets returned.
type Producer interface {
	Get() (*Transaction, error)
}

// Consumer any sink will implement this to consume transactions from the source one at a time.
type Consumer interface {
	Put(*Transaction) error
}

// Transaction is the main Domain model for this product.
// It will contain all the important information of a money transfer transaction for the given system.
type Transaction struct {
	ID     int
	Date   time.Time
	Amount float64
}

// Process will consume a source and aggregate the transactions into the given consumers
// until the source is empty.
func Process(source Producer, aggregators ...Consumer) error {
	for {
		t, err := source.Get()
		if err != nil {
			if errors.Is(err, NoRecordError) {
				return nil
			}
			return err
		}
		for _, a := range aggregators {
			if err := a.Put(t); err != nil {
				return err
			}
		}
	}
}
