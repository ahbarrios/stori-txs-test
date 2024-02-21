package main

import (
	"crypto/tls"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/email"
	"github.com/ahbarrios/stori-txn-test/internal/email/template"
	"github.com/ahbarrios/stori-txn-test/internal/source/csv"
	"github.com/ahbarrios/stori-txn-test/internal/tx"
	"github.com/ahbarrios/stori-txn-test/internal/tx/agg"

	"github.com/emersion/go-smtp"
)

var localSMTP *smtp.Server

func main() {
	// Start the SMTP server
	go func() {
		if err := localSMTP.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	defer localSMTP.Close()

	// Open the file source & handle the error and close the file at the end
	fd, err := os.Open(configmgr.GetSourcePath())
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	// The Source will implement some interface that always be delivered by the consumer (the func below) to
	// read from a file source (csv) and get transactions until io.EOF.
	src := csv.NewHandler(fd)

	// Write a func that receive a Source: it will give Transactions as an output and visit that source for every accumulator.
	// It will have several Accumulators as input args that will be used as visitors for every processed transaction
	// extracted from the Source
	// The return type will be an error
	var (
		balance    agg.Balance
		avgDebit   agg.AvgDebit
		avgCredit  agg.AvgCredit
		txnByMonth = make(agg.MonthlySummary)
	)
	if err := tx.Process(src, &balance, &avgDebit, &avgCredit, txnByMonth); err != nil {
		log.Fatal(err)
	}

	rcp := configmgr.GetRecipient()
	// Write a func that receives an INPUT and return an HTML email template as result
	bodyEmail := func() (io.Reader, error) {
		r, err := template.NewSummaryBody(balance, avgCredit, avgDebit, txnByMonth)
		if err != nil {
			return nil, err
		}

		// Read all bytes from body email
		br, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		// Compound the email body with RFC 822 compatible format
		m := strings.NewReader("To: " + rcp + "\r\n" +
			"Subject: Transaction Summary Email\r\n" +
			"\r\n" +
			string(br) + "\r\n")
		return m, nil
	}

	// work on the SMTP server to send the email
	if err := email.NewClient(serverAddr, nil).SendMail("stori@storicard.com", []string{rcp}, bodyEmail); err != nil {
		log.Fatal(err)
	}
}

func init() {
	s := smtp.NewServer(&email.Backend{})
	s.Addr = serverAddr
	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50

	// enable TLS support for secure clients with self-signed certificate
	cert, err := tls.LoadX509KeyPair("internal/email/testdata/server.pem", "internal/email/testdata/server.key")
	if err != nil {
		log.Fatal(err)
	}
	s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	localSMTP = s
}
