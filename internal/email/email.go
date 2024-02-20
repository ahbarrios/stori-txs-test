// Email will handle the email sending for the given system along with all the supported templates.
// This package will abstract the underlining email client to keep the system decoupled from whatever package that we use to send emails.
package email

import (
	"io"
	"strings"

	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"
)

// SummaryEmailBodyFunc it will create a email body with the given [tx/agg] as input
// that will be used to create the summary email by any template.
type SummaryEmailBodyFunc func(agg.Balance, agg.AvgCredit, agg.AvgDebit, agg.MonthlySummary) (*strings.Reader, error)

// Sender it will represent the email client to send [tx.Transaction] related emails
// The simple implementation will receive the sender, the receiver and the body function of the email that will return [io.Reader] compliant with [RFC 822].
//
// [RFC 822]: https://www.rfc-editor.org/rfc/rfc822.html
type Sender interface {
	SendMail(sender string, recipients []string, bodyFn func() (io.Reader, error)) error
}
