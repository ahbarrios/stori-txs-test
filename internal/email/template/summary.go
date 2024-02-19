package template

import (
	"strings"

	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"
)

func NewSummaryBody(agg.Balance, agg.AvgCredit, agg.AvgDebit, agg.MonthlySummary) (*strings.Reader, error) {
	panic("not implemented")
}
