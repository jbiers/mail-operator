package provider

type EmailData struct {
	ApiToken  string
	Text      string
	Subject   string
	Recipient string
	Sender    string
}

type EmailProvider interface {
	SendEmail(data *EmailData) (error, string, string)
}
