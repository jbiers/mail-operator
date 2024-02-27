package provider

import (
	"context"
	"time"

	"github.com/mailersend/mailersend-go"
)

type MailerSend struct {
}

func (m *MailerSend) SendEmail(data *EmailData) (error, string, string) {
	ms := mailersend.NewMailersend(data.ApiToken)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Email: data.Sender,
	}

	recipients := []mailersend.Recipient{
		{
			Email: data.Recipient,
		},
	}

	message := mailersend.Message{
		From:       from,
		Recipients: recipients,
		Subject:    data.Subject,
		Text:       data.Text,
	}

	res, err := ms.Email.Send(ctx, &message)

	return err, res.Header.Get("x-message-id"), res.Status
}
