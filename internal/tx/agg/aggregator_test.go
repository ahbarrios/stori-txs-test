package agg

import (
	"testing"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/tx"
	"github.com/google/go-cmp/cmp"
)

var data = []tx.Transaction{
	{Date: time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC), Amount: 60.5},
	{Date: time.Date(0, time.July, 28, 0, 0, 0, 0, time.UTC), Amount: -10.3},
	{Date: time.Date(0, time.August, 2, 0, 0, 0, 0, time.UTC), Amount: -20.46},
	{Date: time.Date(0, time.August, 13, 0, 0, 0, 0, time.UTC), Amount: 10},
}

func TestBalance_Put(t *testing.T) {
	var b Balance
	for _, tx := range data {
		if err := b.Put(tx); err != nil {
			t.Errorf("Balance.Put() error = %v", err)
		}
	}
	if b != Balance(39.74) {
		t.Errorf("Balance = %v, want %v", b, 39.74)
	}
}

func TestAvgDebit_Put(t *testing.T) {
	var ad AvgDebit
	for _, tx := range data {
		if err := ad.Put(tx); err != nil {
			t.Errorf("AvgDebit.Put() error = %v", err)
		}
	}
	if ad.Value != 35.25 {
		t.Errorf("AvgDebit = %v, want %v", ad.Value, 35.25)
	}
}

func TestAvgCredit_Put(t *testing.T) {
	var ac AvgCredit
	for _, tx := range data {
		if err := ac.Put(tx); err != nil {
			t.Errorf("AvgCredit.Put() error = %v", err)
		}
	}
	if ac.Value != -15.38 {
		t.Errorf("AvgCredit = %v, want %v", ac.Value, -15.38)
	}
}

func TestMonthlySummary_Put(t *testing.T) {
	ms := make(MonthlySummary)
	for _, tx := range data {
		if err := ms.Put(tx); err != nil {
			t.Errorf("MonthlySummary.Put() error = %v", err)
		}
	}
	tg := map[time.Month]int{
		time.July:   2,
		time.August:   2,
	}
	if !cmp.Equal(ms, MonthlySummary(tg)) {
		t.Errorf("MonthlySummary = %v, want %v", ms, tg)
	}
}
