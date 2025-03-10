package email

import (
	"app/pkg/config"
	"app/pkg/mail"
	"log"
	"strconv"
)

var (
	smtpClient *mail.SMTPClient
)

func MustLoad() {
	smtpHost := config.ENV().MAIL_SMTP_HOST
	smtpPort, err := strconv.Atoi(config.ENV().MAIL_SMTP_PORT)
	if err != nil {
		log.Fatalf("Error converting SMTP port to int: %v", err)
	}
	smtpUser := config.ENV().MAIL_SMTP_USERNAME
	smtpPassword := config.ENV().MAIL_SMTP_PASSWORD
	smtpTLS := config.ENV().MAIL_SMTP_TLS
	smtpClient = mail.GetSMTPClient(smtpHost, smtpPort, smtpUser, smtpPassword, smtpTLS)
	log.Printf("✅ Successfully connected to SMTP %s:%d", smtpHost, smtpPort)
}

func SendEmail(from, to, subject, body string) error {
	return smtpClient.SendEmail(from, to, subject, body)
}
