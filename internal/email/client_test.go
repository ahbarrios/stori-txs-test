package email

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/emersion/go-sasl"
)

var csr = x509.CertificateRequest{
	Subject: pkix.Name{
		Country:            []string{"MX"},
		Organization:       []string{"Stori"},
		OrganizationalUnit: []string{"Security"},
		Locality:           []string{"CDMX"},
		CommonName:         "storicard.com",
	},
	DNSNames:    []string{"*.storicard.com", "localhost"},
	IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
}

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
