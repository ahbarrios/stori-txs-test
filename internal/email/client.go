package email

import (
	"io"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Client will implement the [email.Sender] interface to send emails using the [smtp.Client] from the [smtp] package.
type Client struct {
	addr string
	auth sasl.Client
}

// SendMail will send an email using [smtp.SendMail] method.
func (c *Client) SendMail(sender string, recipients []string, bodyFn func() (io.Reader, error)) error {
	body, err := bodyFn()
	if err != nil {
		return err
	}
	return smtp.SendMail(c.addr, c.auth, sender, recipients, body)
}

// NewClient will create a new [email.Client] with the given address and authentication method.
func NewClient(addr string, auth sasl.Client) *Client {
	return &Client{addr: addr, auth: auth}
}
