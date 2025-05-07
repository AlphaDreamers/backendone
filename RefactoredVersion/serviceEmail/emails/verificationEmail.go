package emails

import "github.com/resend/resend-go/v2"

var From = "onboarding@resend.dev"

func VerificationEmail(to string, link string) *resend.SendEmailRequest {
	return &resend.SendEmailRequest{
		From:    From,
		To:      []string{to},
		Subject: "Verify Your Email Address",
		Html: `
		<h2>Verify Your Email</h2>
		<p>Please click the button below to verify your email:</p>
		<a href="` + link + `" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none;">Verify Email</a>
		<p>If you did not request this, you can ignore this email.</p>
		`,
	}
}

func PasswordResetEmail(to string, link string) *resend.SendEmailRequest {
	return &resend.SendEmailRequest{
		From:    From,
		To:      []string{to},
		Subject: "Reset Your Password",
		Html: `
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password:</p>
		<a href="` + link + `" style="background-color: #FF5722; color: white; padding: 10px 20px; text-decoration: none;">Reset Password</a>
		<p>This link expires in 24 hours. If you didn’t request it, ignore this email.</p>
		`,
	}
}

func OrderConfirmationEmail(to string, orderID string, items []string) *resend.SendEmailRequest {
	return &resend.SendEmailRequest{
		From:    "orders@yourdomain.com",
		To:      []string{to},
		Subject: "Order Confirmation - #" + orderID,
		Html: `
		<h2>Your Order Has Been Received</h2>
		<p>Order ID: ` + orderID + `</p>
		<ul>` + generateOrderItemsHTML(items) + `</ul>
		<p>Thank you for your purchase!</p>
		`,
	}
}

func generateOrderItemsHTML(items []string) string {
	html := ""
	for _, item := range items {
		html += "<li>" + item + "</li>"
	}
	return html
}
