package source

import (
	"bytes"
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/tx"
	"github.com/google/go-cmp/cmp"
)

//go:embed files/csv_test.csv
var testOKSource []byte

//go:embed files/empty_csv_test.csv
var testNoRecordSource []byte

//go:embed files/no_csv.txt
var testBadSource []byte

//go:embed files/empty_file.txt
var testEmptySource []byte

func TestCSV_Read(t *testing.T) {
	c := NewHandler(bytes.NewReader(testOKSource))
	tests := []struct {
		name    string
		want    *tx.Transaction
		err     error
		wantErr bool
	}{
		{
			name: "Skip the header",
			want: &tx.Transaction{
				Date:   time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC),
				Amount: 60.5,
			},
		},
		{
			name: "Read 2nd row",
			want: &tx.Transaction{
				ID:     1,
				Date:   time.Date(0, time.July, 28, 0, 0, 0, 0, time.UTC),
				Amount: -10.3,
			},
		},
		{
			name: "Read 3rd row",
			want: &tx.Transaction{
				ID:     2,
				Date:   time.Date(0, time.August, 2, 0, 0, 0, 0, time.UTC),
				Amount: -20.46,
			},
		},
		{
			name: "Read last row",
			want: &tx.Transaction{
				ID:     3,
				Date:   time.Date(0, time.August, 13, 0, 0, 0, 0, time.UTC),
				Amount: 10,
			},
		},
		{
			name:    "Read after EOF",
			err:     tx.NoRecordError,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("CSV.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.err {
				t.Errorf("CSV.Read() error = %v, want %v", err, tt.err)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("CSV.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSV_NoRecordSource(t *testing.T) {
	c := NewHandler(bytes.NewReader(testNoRecordSource))

	_, err := c.Get()
	if err != tx.NoRecordError {
		t.Errorf("CSV.Read() error = %v, expecting NoRecordError", err)
	}
}

func TestCSV_NoCSVSource(t *testing.T) {
	c := NewHandler(bytes.NewReader(testBadSource))

	_, err := c.Get()
	if !errors.Is(err, BadRecordError{}) {
		t.Errorf("CSV.Read() error = %v, expecting BadRecordError", err)
	}
}

func TestCSV_EmptySource(t *testing.T) {
	c := NewHandler(bytes.NewReader(testEmptySource))

	_, err := c.Get()
	if err != tx.NoRecordError {
		t.Errorf("CSV.Read() error = %v, expecting NoRecordError", err)
	}
}
