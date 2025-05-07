package utils

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"github.com/spf13/viper"
)

func SendEmail(tokenCode, email string) {
	apiKey := viper.GetString("resend.api.key")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Html:    "<strong>" + tokenCode + "</strong>",
		Subject: "Hello from Golang",
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sent.Id)
}
