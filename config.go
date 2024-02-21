package main

const (
	serverAddr = "localhost:1025"
	recipient  = "adrian2monk@gmail.com"
	sourceFile = "examples/txns.csv"
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
	return serverAddr
}

func (c *config) GetRecipient() string {
	return recipient
}

func (c *config) GetSourcePath() string {
	return sourceFile
}
