package email

import (
	"app/pkg/config"
	"app/pkg/mail"
	"errors"
	"log"
	"strconv"
	"sync"
)

var (
	smtpClient *mail.SMTPClient
	once       sync.Once

	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	smtpTLS      bool
)

func Mail() {
	once.Do(func() {
		smtpHost = config.ENV().MAIL_SMTP_HOST
		smtpPort, err := strconv.Atoi(config.ENV().MAIL_SMTP_PORT)
		if err != nil {
			log.Fatalf("Error converting SMTP port to int: %v", err)
		}
		smtpUser = config.ENV().MAIL_SMTP_USERNAME
		smtpPassword = config.ENV().MAIL_SMTP_PASSWORD
		smtpTLS = config.ENV().MAIL_SMTP_TLS
		smtpClient = mail.GetSMTPClient(smtpHost, smtpPort, smtpUser, smtpPassword, smtpTLS)
		log.Printf("✅ Successfully connected to SMTP %s:%d", smtpHost, smtpPort)
	})
}

func SendEmail(to, subject, body string) error {
	if smtpClient == nil {
		return errors.New("SMTP client not initialized")
	}
	return smtpClient.SendEmail(smtpUser, to, subject, body)
}
