package email

import (
	"crypto/tls"
	"log"
	"testing"
	"time"

	"github.com/emersion/go-smtp"
)

const testServerAddr = "localhost:1025"

var testServer *smtp.Server

func TestMain(m *testing.M) {
	s := smtp.NewServer(&Backend{})
	s.Addr = testServerAddr
	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50

	// enable TLS support for secure clients with self-signed certificate
	cert, err := tls.LoadX509KeyPair("testdata/server.pem", "testdata/server.key")
	if err != nil {
		log.Fatal(err)
	}
	s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	testServer = s
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	m.Run()
	s.Close()
}

func TestServer(t *testing.T) {
	// Connect to the remote SMTP server.
	_, err := smtp.Dial(testServerAddr)
	if err != nil {
		log.Fatal(err)
	}

}
