package sender

import (
	"github.com/resend/resend-go/v2"
)

type EmailContent struct {
	From    string
	To      []string
	Subject string
	HTML    string
}

func SendEmail(client *resend.Client, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	return client.Emails.Send(params)
}
