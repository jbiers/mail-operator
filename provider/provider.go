package provider

import (
	"context"
	"time"

	mailersend "github.com/mailersend/mailersend-go"
)

//I think this file should not be a separate package but actually a part of controllar. but lets see

type EmailData struct {
	ApiToken  string
	Text      string
	Subject   string
	Recipient string
	Sender    string
}

func SendEmail(email *EmailData) (error, string, string) {
	ms := mailersend.NewMailersend(email.ApiToken)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Email: email.Sender,
	}

	recipients := []mailersend.Recipient{
		{
			Email: email.Recipient,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(email.Subject)
	message.SetText(email.Text)

	res, err := ms.Email.Send(ctx, message)

	return err, res.Header.Get("x-message-id"), res.Status
}
