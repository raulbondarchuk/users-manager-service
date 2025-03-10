package mail

import (
	"crypto/tls"
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	dialer *gomail.Dialer
}

var (
	instance *SMTPClient
	once     sync.Once
)

// GetSMTPClient returns the single instance of SMTPClient with TLS support
func GetSMTPClient(smtpHost string, smtpPort int, username, password string, hasTLS bool) *SMTPClient {
	once.Do(func() {
		instance = &SMTPClient{
			dialer: &gomail.Dialer{
				Host:     smtpHost,
				Port:     smtpPort,
				Username: username,
				Password: password,
				SSL:      false, // SSL false, para utilizar TLS
				TLSConfig: &tls.Config{
					InsecureSkipVerify: hasTLS, // TLS
				},
			},
		}
	})
	return instance
}

func (c *SMTPClient) SendEmail(from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send email
	if err := c.dialer.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}
