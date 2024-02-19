// Aggregator contain all the necessary accumulators to extract insight and give final users
// summarized information about its finantial transactions within the system.
package agg

import (
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/tx"
)

// average it will represent the average of a given set of values
// since we need two values to calculate the average we will use this internal type
type average struct {
	n     int
	Value float64
}

func (a *average) Add(v float64) {
	a.n++
	a.Value += (v - a.Value) / float64(a.n)
}

// Balance total balance of the source
type Balance float64

// Put it will implement [internal/tx/Consumer] and process a Transaction value as input
// to produce the total balance of the source
func (b *Balance) Put(t *tx.Transaction) error {
	*b += Balance(t.Amount)
	return nil
}

// AvgCredit average debit transaction amount
type AvgCredit struct {
	average
}

// Put it will implement [internal/tx/Consumer] and process a Transaction value as input
// to produce the average debit transaction amount. A debit transaction will be identified by
// a non-negative amount.
func (ad *AvgCredit) Put(t *tx.Transaction) error {
	if t.Amount >= 0 {
		ad.average.Add(t.Amount)
	}
	return nil
}

// AvgDebit average credit transaction amount
type AvgDebit struct {
	average
}

// Put it will implement [internal/tx/Consumer] and process a Transaction value as input
// to produce the average credit transaction amount. A credit transaction will be identified by
// a negative amount.
func (ac *AvgDebit) Put(t *tx.Transaction) error {
	if t.Amount < 0 {
		ac.average.Add(t.Amount)
	}
	return nil
}

// MonthlySummary number of transactions by month
type MonthlySummary map[time.Month]int

// Put it will implement [internal/tx/Consumer] and process a Transaction value as input
// to produce the number of transactions by [time.Month].
func (ms MonthlySummary) Put(t *tx.Transaction) error {
	if _, ok := ms[t.Date.Month()]; !ok {
		ms[t.Date.Month()] = 0
	}
	ms[t.Date.Month()]++
	return nil
}
