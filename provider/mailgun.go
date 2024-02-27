package provider

import (
	"context"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailGun struct {
}

func (m *MailGun) SendEmail(data *EmailData) (error, string, string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	domain := strings.Split(data.Sender, "@")[1]

	mg := mailgun.NewMailgun(domain, data.ApiToken)
	message := mg.NewMessage(
		data.Sender,
		data.Subject,
		data.Text,
		data.Recipient,
	)

	status, messageId, err := mg.Send(ctx, message)

	return err, messageId, status
}
