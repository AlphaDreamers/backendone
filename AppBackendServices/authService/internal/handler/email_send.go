package handler

import (
	"github.com/goccy/go-json"
	"github.com/streadway/amqp"
)

func (i *Impl) SendEmail(userName, email, code string) {
	i.logger.Infof("Sending email to %s", email)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		i.logger.Errorf("Failed to connect to RabbitMQ")
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			i.logger.Errorf("Failed to close connection")
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		i.logger.Errorf("Failed to open a channel")
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			i.logger.Errorf("Failed to close channel")
		}
	}(ch)

	q, err := ch.QueueDeclare("email",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		i.logger.Errorf("Failed to declare a queue")
	}

	type send struct {
		Email string `json:"email"`
		User  string `json:"user_name"`
		Code  string `json:"code"`
	}

	body := send{
		Email: email,
		User:  userName,
		Code:  code,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		i.logger.Errorf("Failed to marshal body")
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})

	if err != nil {
		i.logger.Errorf("Failed to publish a message")
	}
	i.logger.Infof("Sent email to %s", email)
}
