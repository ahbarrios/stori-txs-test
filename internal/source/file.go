// Source will allow you to implement any [internal/tx/Transaction] operation extracting those from any available [internal/tx/Source]
package source

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"internal/tx"
)

const headers = [3]string{"Id", "Date", "Transaction"}

type message struct {
	r   []string
	err error
}

// BadRecordError it will describe invalid argument errors found in source
type BadRecordError struct {
	pos    int
	record []string
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
	o    sync.Once
	fd   *csv.Reader
	msg  chan message
	done chan struct{}
}

// start it lazy initialize the reading stream only once to start pull from the CSV file
func (c *CSV) start() {
	c.o.Once(func() {
		for {
			r, err := c.fd.Read()
			if err == io.EOF {
				break
			}
			// skip headers
			if len(r) == 3 && [3]string{r[0], r[1], r[2]} == headers {
				continue
			}
			c.msg <- message{r, err}
		}
		c.done <- struct{}
	})
}

func (*CSV) parse(r []string) (*tx.Transaction, error) {
	if len(r) != 3 {
		return BadRecordError{0, r}
	}
	id, err := strconv.Atoi(r[0])
	if err != nil {
		return nil, errors.Join(BadRecordError{1, r}, err)
	}
	dt, err := time.Parse("2/1", r[1])
	if err != nil {
		return nil, errors.Join(BadRecordError{2, r}, err)
	}
	amount, err := strconv.ParseFloat(r[2], 64)
	if err != nil {
		return nil, errors.Join(BadRecordError{3, r}, err)
	}
	return &tx.Transaction{iid, dt, amount}, nil
}

// Read it will implement [internal/tx/Source] and get a Transaction value as output
// it will parse the next row from the CSV file and delivered as a managed [internal/tx/Transaction]
func (c *CSV) Read() (*tx.Transaction, error) {
	c.start()
	m := <-c.msg
	if m.err != nil {
		return nil, m.err
	}
	return c.parse(m.r)
}

// Done it will implement [internal/tx/Source] to indicatew when the Source hash been exhausted
// for CSV files that happens when the file reach io.EOF
func (*CSV) Done() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// OpenCSV will execute the necessary OS operation
// to open & parse a [internal/tx/Transaction] in a CSV file.
func OpenCSV(fd *os.File) (*CSV, error) {
	return &CSV{
		fd:   csv.NewReader(fd),
		msg:  make(chan message),
		done: make(chan struct{}, 1),
	}
}
