package main

import (
	"os"

	"github.com/emersion/go-sasl"
)

const localSMTPAddr = "localhost:1025"

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

func (c *config) GetSender() string {
	return os.Getenv("STORI_SENDER")
}

func (c *config) GetRecipient() string {
	return os.Getenv("STORI_RECIPIENT")
}

func (c *config) GetSourcePath() string {
	return os.Getenv("STORI_SOURCE_PATH")
}

func (c *config) GetAuth() sasl.Client {
	if c.IsLocal() {
		return sasl.NewPlainClient("", "username", "password")
	}

	usr := os.Getenv("STORI_SMTP_USERNAME")
	pwd := os.Getenv("STORI_SMTP_PASSWORD")
	return sasl.NewPlainClient("", usr, pwd)
}

// IsLocal will return true if the current environment is local to spin up the mock SMTP server
func (c *config) IsLocal() bool {
	return c.GetSMTPServer() == localSMTPAddr
}
