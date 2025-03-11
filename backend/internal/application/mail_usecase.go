package application

import (
	"app/internal/infrastructure/transport/email"
	"errors"
	"strings"
)

type MailUseCase struct {
}

func NewMailUseCase() *MailUseCase {
	return &MailUseCase{}
}

func (uc *MailUseCase) SendEmailForgotPassword(to, subject, body, link string) error {

	// Check if body contains {link}
	if !strings.Contains(body, "{link}") {
		return errors.New("incorrect email format: missing {link}")
	}

	// Replace {link} with link
	body = strings.ReplaceAll(body, "{link}", link)
	if strings.Contains(body, "{username}") {
		body = strings.ReplaceAll(body, "{username}", to)
	}

	// Send the email
	return email.SendEmail(to, subject, body)
}
