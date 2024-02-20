package email

import (
	"io"
	"strings"
	"testing"

	"github.com/ahbarrios/stori-txn-test/internal/email/template"
	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"
	"github.com/emersion/go-sasl"
)

func TestClient_SendMail(t *testing.T) {
	a := sasl.NewPlainClient("", "username", "password")
	c := NewClient(testServerAddr, a)

	m := strings.NewReader("To: to@example.com\r\n" +
		"Subject: testing Client.SendMail!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := c.SendMail("sender@localhost", []string{"to@example.com"}, func() (io.Reader, error) {
		return m, nil
	})
	if err != nil {
		t.Fatalf("SendMail() error = %v", err)
	}
}

func TestClient_SendMail_Summary(t *testing.T) {
	a := sasl.NewPlainClient("", "username", "password")
	c := NewClient(testServerAddr, a)

	err := c.SendMail("sender@localhost", []string{"to@example.com"}, func() (io.Reader, error) {
		var (
			c agg.AvgCredit
			d agg.AvgDebit
		)
		r, err := template.NewSummaryBody(agg.Balance(39.6), c, d, nil)
		if err != nil {
			return nil, err
		}

		// Read all bytes from body email
		br, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		// Compound the email body with RFC 822 compatible format
		m := strings.NewReader("To: to@example.com\r\n" +
			"Subject: testing Client.SendMail Summary Email!\r\n" +
			"\r\n" +
			string(br) + "\r\n")
		return m, err
	})
	if err != nil {
		t.Fatalf("SendMail() error = %v", err)
	}
}
