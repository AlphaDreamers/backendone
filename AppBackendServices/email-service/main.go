package main

import (
	"bytes"
	"github.com/go-mail/mail"
	"github.com/goccy/go-json"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"net/http"

	"html/template"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type EmailContent struct {
	Name string
	Code string
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

func SendWelcomeEmail(recipient string, data EmailContent) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	temfile, err := template.ParseFiles("./email.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := temfile.Execute(&body, data); err != nil {
		return err
	}
	m := mail.NewMessage()
	m.SetHeader("From", smtpFrom)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", data.Name)
	m.SetBody("text/html", body.String())
	m.SetHeader("Content-Type", "text/html; charset=utf-8")
	m.SetHeaders(map[string][]string{
		"X-Auto-Response-Suppress": {"All"},
		"Auto-Submitted":           {"auto-generated"},
	})
	dialer := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func consume() {

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

	log.Printf("Connecting to RabbitMQ")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(conn)

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer channel.Close()
	q, err := channel.QueueDeclare(
		"email",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	type send struct {
		Email string `json:"email"`
		User  string `json:"user_name"`
		Code  string `json:"code"`
	}
	var body send

	for msg := range msgs {
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Fatal(err)
		}
		log.Printf("Received a message: %s", body.Email)
		go func() {
			err := SendWelcomeEmail(body.Email, EmailContent{
				Name: body.User,
				Code: body.Code,
			})
			if err != nil {
				log.Printf(err.Error())
			}
		}()
	}
}

//docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ConsulRegister()
	go consume()
	log.Println("Starting server")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func ConsulRegister() {
	config := api.DefaultConfig()
	config.Address = "http://consul:8500"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	serviceRegistration := &api.AgentServiceRegistration{
		ID:      "email-service",
		Name:    "email",
		Address: "email-service:8080",
		Port:    8080,
		Tags:    []string{"email"},
		Check: &api.AgentServiceCheck{
			HTTP:     "http://emai-service:9008/health",
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	if err := client.Agent().ServiceRegister(serviceRegistration); err != nil {
		log.Fatal(err.Error())
	}

}
