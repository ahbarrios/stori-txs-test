package template

import (
	"bytes"
	_ "embed"
	"io"
	"testing"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/tx"
	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"
	"github.com/google/go-cmp/cmp"
)

//go:embed summary_empty_test.html
var emptyTest []byte

//go:embed summary_example_test.html
var validTest []byte

var data = []tx.Transaction{
	{Date: time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC), Amount: 60.5},
	{Date: time.Date(0, time.July, 28, 0, 0, 0, 0, time.UTC), Amount: -10.3},
	{Date: time.Date(0, time.August, 2, 0, 0, 0, 0, time.UTC), Amount: -20.46},
	{Date: time.Date(0, time.August, 13, 0, 0, 0, 0, time.UTC), Amount: 10},
}

func TestNewSummaryBody_Empty(t *testing.T) {
	var (
		b agg.Balance
		c agg.AvgCredit
		d agg.AvgDebit
	)
	r, err := NewSummaryBody(b, c, d, nil)
	if err != nil {
		t.Fatalf("NewSummaryBody() error = %v", err)
	}

	br, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	if !cmp.Equal(br, emptyTest) {
		t.Errorf("NewSummaryBody() = %v, want %s", r, emptyTest)
		t.Logf("Diff: %s", cmp.Diff(br, emptyTest))
	}
}

func TestNewSummaryBody_Valid(t *testing.T) {
	b := agg.Balance(39.74)
	c := newAvgCredit(t)
	d := newAvgDebit(t)
	m := newMonthlySummary(t)
	r, err := NewSummaryBody(b, c, d, m)
	if err != nil {
		t.Fatalf("NewSummaryBody() error = %v", err)
	}

	br, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	// remove trailing spaces out of output to match the saved/formatted file
	br = bytes.ReplaceAll(br, []byte("                                \n"), []byte(""))
	if !cmp.Equal(br, validTest) {
		t.Errorf("NewSummaryBody() = %s, want %s", br, validTest)
		t.Logf("Diff: %s", cmp.Diff(br, validTest))
	}
}

func newAvgCredit(t *testing.T) agg.AvgCredit {
	t.Helper()
	var ad agg.AvgCredit
	for _, tx := range data {
		if err := ad.Put(&tx); err != nil {
			t.Errorf("AvgCredit.Put() error = %v", err)
		}
	}
	return ad
}

func newAvgDebit(t *testing.T) agg.AvgDebit {
	t.Helper()
	var ac agg.AvgDebit
	for _, tx := range data {
		if err := ac.Put(&tx); err != nil {
			t.Errorf("AvgDebit.Put() error = %v", err)
		}
	}
	return ac
}

func newMonthlySummary(t *testing.T) agg.MonthlySummary {
	t.Helper()
	ms := make(agg.MonthlySummary)
	for _, tx := range data {
		if err := ms.Put(&tx); err != nil {
			t.Errorf("MonthlySummary.Put() error = %v", err)
		}
	}
	return ms
}
