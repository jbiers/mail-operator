package controller

import "github.com/jbiers/mail-operator/internal/provider"

type EmailSender struct {
	provider provider.EmailProvider
}

func initEmailSender(p provider.EmailProvider) *EmailSender {
	return &EmailSender{
		provider: p,
	}
}

func (e *EmailSender) sendEmail(d *provider.EmailData) (error, string, string) {
	return e.provider.SendEmail(d)
}
