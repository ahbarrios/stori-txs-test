package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/ahbarrios/stori-txn-test/internal/email"
	"github.com/emersion/go-smtp"
)

// configmgr will be a global variable that will be used to get the configuration from the environment
var configmgr config

// config will be the configuration manager for the given system
// this could be a package that will be used to get the configuration from the environment
// following the 12-factor app methodology.
//
// The current implementation only suuport environment variables but it could be extended to support
// other sources like yaml, databases, configmaps, cli arguments, etc.
//
// [12-factor app]: https://12factor.net/config
type config struct{}

func (c *config) GetSMTPServer() string {
	return os.Getenv("STORI_SMTP_SERVER")
}

func (c *config) GetRecipient() string {
	return os.Getenv("STORI_RECIPIENT")
}

func (c *config) GetSourcePath() string {
	return os.Getenv("STORI_SOURCE_PATH")
}

// IsLocal will return true if the current environment is local to spin up the mock SMTP server
func (c *config) IsLocal() bool {
	return c.GetSMTPServer() == "localhost:1025"
}

func init() {
	os.Setenv("STORI_SMTP_SERVER", "localhost:1025")
	os.Setenv("STORI_RECIPIENT", "adrian2monk@gmail.com")
	os.Setenv("STORI_SOURCE_PATH", "examples/txns.csv")

	if configmgr.IsLocal() {
		s := smtp.NewServer(&email.Backend{})
		s.Addr = configmgr.GetSMTPServer()
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
}
