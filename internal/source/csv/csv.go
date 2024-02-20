// CSV will allow you to implement any [internal/tx/Transaction] operation extracting those from any available [internal/tx/Source]
package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/tx"
)

var headers = [3]string{"Id", "Date", "Transaction"}

// BadRecordError it will describe invalid argument errors found in source
type BadRecordError struct {
	pos    int
	record []string
}

func (e BadRecordError) Is(target error) bool {
	_, ok := target.(BadRecordError)
	return ok
}

func (e BadRecordError) Error() string {
	if len(e.record) != 3 {
		return fmt.Sprintf("wrong number of elements found in row, want 3 got %d", len(e.record))
	}
	return fmt.Sprintf("invalid value found %v in column %d", e.record, e.pos)
}

// CSV it will represent the CSV source file for transactions with
// all the operations to handle I/O available
type CSV struct {
	fd *csv.Reader
}

// Get it will implement [tx.Producer] and get a Transaction value as output
// it will parse the next row from the CSV file and parse as the managed [internal/tx/Transaction]
func (c *CSV) Get() (*tx.Transaction, error) {
	r, err := c.read()
	if err != nil {
		return nil, err
	}
	return c.parse(r)
}

// read it will consume the reading stream from the CSV file
// and return [tx.NoRecordError] when the file is empty
func (c *CSV) read() ([]string, error) {
	read := func(f *csv.Reader) ([]string, error) {
		r, err := f.Read()
		if err == io.EOF {
			return nil, tx.NoRecordError
		}
		return r, err
	}

	r, err := read(c.fd)
	// skip headers
	if len(r) == 3 && [3]string{r[0], r[1], r[2]} == headers {
		r, err = read(c.fd)
	}
	return r, err
}

func (*CSV) parse(r []string) (*tx.Transaction, error) {
	if len(r) != 3 {
		return nil, BadRecordError{0, r}
	}
	id, err := strconv.Atoi(r[0])
	if err != nil {
		return nil, errors.Join(BadRecordError{1, r}, err)
	}
	dt, err := time.Parse("1/2", r[1])
	if err != nil {
		return nil, errors.Join(BadRecordError{2, r}, err)
	}
	amount, err := strconv.ParseFloat(r[2], 64)
	if err != nil {
		return nil, errors.Join(BadRecordError{3, r}, err)
	}
	return &tx.Transaction{
		ID:     id,
		Date:   dt,
		Amount: amount,
	}, nil
}

// NewHandler will execute the necessary OS operation to parse a [tx.Transaction] in a CSV file.
func NewHandler(fd io.Reader) *CSV {
	return &CSV{csv.NewReader(fd)}
}
