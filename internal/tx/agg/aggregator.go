// Aggregator contain all the necessary accumulators to extract insight and give final users
// summarized information about its finantial transactions within the system.
package agg

import "time"

// Balance total balance of the source
type Balance float64

// AvgDebit average debit transaction amount
type AvgDebit float64

// AvgCredit average credit transaction amount
type AvgCredit float64

// MonthlySummary number of transactions by month within the current year
type MontlySumary struct {
	Count int
	Month time.Month
}
