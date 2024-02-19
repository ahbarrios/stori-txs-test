package tx

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var data = []Transaction{
	{Date: time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC), Amount: 60.5},
	{Date: time.Date(0, time.July, 28, 0, 0, 0, 0, time.UTC), Amount: -10.3},
	{Date: time.Date(0, time.August, 2, 0, 0, 0, 0, time.UTC), Amount: -20.46},
	{Date: time.Date(0, time.August, 13, 0, 0, 0, 0, time.UTC), Amount: 10},
}

type fakeProducer struct {
	i int
}

func (f *fakeProducer) Get() (*Transaction, error) {
	if f.i >= len(data) {
		return nil, NoRecordError
	}

	t := data[f.i]
	f.i++
	return &t, nil
}

type fakeConsumer struct {
	i int
}

func (f *fakeConsumer) Put(t *Transaction) error {
	if !cmp.Equal(*t, data[f.i]) {
		return fmt.Errorf("Transaction = %v, want %v", *t, data[f.i])
	}
	f.i++
	return nil
}

func TestProcess(t *testing.T) {
	if err := Process(&fakeProducer{}, &fakeConsumer{}); err != nil {
		t.Errorf("Process() error = %v", err)
	}
}
