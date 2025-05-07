package main

import (
	"context"
	"fmt"
	"log"
	"os"

	elasticemail "github.com/elasticemail/elasticemail-go"
)

type Order struct {
	ID          string
	Customer    Customer
	Items       []Item
	TotalAmount float64
	Status      string
}

type Customer struct {
	Name  string
	Email string
}

type Item struct {
	Name     string
	Quantity int
	Price    float64
}

func main() {
	// Initialize the client
	cfg := elasticemail.NewConfiguration()
	client := elasticemail.NewAPIClient(cfg)

	// Get API key from environment variable
	apiKey := os.Getenv("6759B1BC36EBEC4721A96FB167EC18BC6BB0DB8B8697DF18B377321B264AA90FBCD9D84E73757CF726F647F7B75898F2")
	if apiKey == "" {
		log.Fatal("ELASTICEMAIL_API_KEY environment variable not set")
	}

	// Create a sample order
	order := Order{
		ID: "ORD-12345",
		Customer: Customer{
			Name:  "John Doe",
			Email: "swanhtet102002@gmail.com",
		},
		Items: []Item{
			{Name: "Product A", Quantity: 2, Price: 25.99},
			{Name: "Product B", Quantity: 1, Price: 59.99},
		},
		TotalAmount: 111.97,
		Status:      "Processing",
	}

	// Send order confirmation
	err := sendOrderConfirmation(client, apiKey, order)
	if err != nil {
		log.Fatalf("Failed to send order confirmation: %v", err)
	}

	fmt.Println("Order confirmation sent successfully!")
}

func sendOrderConfirmation(client *elasticemail.APIClient, apiKey string, order Order) error {
	// Prepare email content
	subject := fmt.Sprintf("Order Confirmation #%s", order.ID)
	htmlContent := generateOrderConfirmationHTML(order)
	textContent := generateOrderConfirmationText(order)

	// Prepare the email data
	emailData := elasticemail.EmailTransactionalMessageData{
		Recipients: elasticemail.TransactionalRecipient{
			To: []string{order.Customer.Email},
		},
		Content: elasticemail.EmailContent{
			Body: []elasticemail.BodyPart{
				{
					ContentType: "HTML",
					Content:     &htmlContent,
				},
				{
					ContentType: "PlainText",
					Content:     &textContent,
				},
			},
			From:    "swanhtetaungp@gmail.com", // Must be verified in Elastic Email
			Subject: &subject,
		},
	}

	// Create context with API key
	ctx := context.WithValue(context.Background(), elasticemail.ContextAPIKeys, map[string]elasticemail.APIKey{
		"apikey": {Key: apiKey},
	})

	// Send the transactional email
	response, _, err := client.EmailsAPI.EmailsTransactionalPost(ctx).EmailTransactionalMessageData(emailData).Execute()
	if err != nil {
		return fmt.Errorf("API error: %v", err)
	}

	log.Printf("Email sent! Transaction ID: %s", *response.TransactionID)
	return nil
}

func generateOrderConfirmationHTML(order Order) string {
	// Generate HTML content for the order confirmation
	itemsHTML := ""
	for _, item := range order.Items {
		itemsHTML += fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td>%d</td>
				<td>$%.2f</td>
				<td>$%.2f</td>
			</tr>`,
			item.Name, item.Quantity, item.Price, float64(item.Quantity)*item.Price)
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Order Confirmation</title>
	</head>
	<body>
		<h1>Thank you for your order, %s!</h1>
		<p>Your order #%s is currently <strong>%s</strong>.</p>
		
		<h2>Order Summary</h2>
		<table border="1" cellpadding="5" cellspacing="0">
			<thead>
				<tr>
					<th>Product</th>
					<th>Quantity</th>
					<th>Unit Price</th>
					<th>Total</th>
				</tr>
			</thead>
			<tbody>
				%s
			</tbody>
			<tfoot>
				<tr>
					<th colspan="3" align="right">Order Total:</th>
					<th>$%.2f</th>
				</tr>
			</tfoot>
		</table>
		
		<p>We'll notify you when your order ships.</p>
		<p>If you have any questions, please contact our support team.</p>
	</body>
	</html>`,
		order.Customer.Name, order.ID, order.Status, itemsHTML, order.TotalAmount)
}

func generateOrderConfirmationText(order Order) string {
	// Generate plain text content for the order confirmation
	itemsText := ""
	for _, item := range order.Items {
		itemsText += fmt.Sprintf("%s (Qty: %d) - $%.2f each - $%.2f total\n",
			item.Name, item.Quantity, item.Price, float64(item.Quantity)*item.Price)
	}

	return fmt.Sprintf(`
Thank you for your order, %s!

Your order #%s is currently %s.

Order Summary:
%s
Order Total: $%.2f

We'll notify you when your order ships.

If you have any questions, please contact our support team.`,
		order.Customer.Name, order.ID, order.Status, itemsText, order.TotalAmount)
}
