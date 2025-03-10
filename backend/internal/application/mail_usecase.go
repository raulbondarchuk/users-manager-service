package application

import "app/internal/infrastructure/transport/email"

type MailUseCase struct {
}

func NewMailUseCase() *MailUseCase {
	return &MailUseCase{}
}

func (uc *MailUseCase) SendEmail(to, subject, body string) error {
	return email.SendEmail(to, subject, body)
}
