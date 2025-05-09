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

func PasswordRelatedEmail(email, url, token string) {
	apiKey := viper.GetString("resend.api.key")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Reset Your Password",
		Html: `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>Password Reset</title>
			</head>
			<body style="font-family: Arial, sans-serif; line-height: 1.6;">
				<h2>Hello,</h2>
				<p>We received a request to reset your password. Please use the button below to set a new password:</p>
				<p>
					<a href="` + url + `" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px;">
						Reset Password
					</a>
				</p>
				<p>If the button doesn't work, copy and paste the following URL into your browser:</p>
				<p><code>` + token + `</code></p>
				<p>This link will expire in 15 minutes.</p>
				<p>If you didn't request a password reset, you can safely ignore this email.</p>
				<br>
				<p>Thanks,<br>The Acme Team</p>
			</body>
			</html>
		`,
		ReplyTo: "no-reply@acme.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return
	}
	fmt.Println("Password reset email sent with ID:", sent.Id)
}
