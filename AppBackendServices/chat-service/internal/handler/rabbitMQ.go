package handler

import (
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
)

var MessageForOfflineUser *amqp.Channel

// RabbitMqInit initializes RabbitMQ connection and channel
func RabbitMqInit() error {
	loggger := logutil.GetLogger()

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		loggger.WithError(err).Warn("Failed to connect to RabbitMQ")
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		loggger.WithError(err).Warn("Failed to open RabbitMQ channel")
		return err
	}

	err = channel.ExchangeDeclare("chat_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		loggger.WithError(err).Warn("Failed to declare RabbitMQ exchange")
		return err
	}
	_, err = channel.QueueDeclare("user_queues", true, false, false, false, nil)
	if err != nil {
		loggger.WithError(err).Warn("Failed to declare user queues")
		return err
	}
	MessageForOfflineUser = channel

	return nil
}

// PublishingToQueue publishes a message to the RabbitMQ exchange
func PublishingToQueue(msg *model.ChatMessage, logger *logrus.Logger) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.WithError(err).Warn("Failed to marshal chat message to JSON")
		return
	}

	err = MessageForOfflineUser.Publish("chat_exchange", msg.RecipientID.String(), false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBytes,
		DeliveryMode: amqp.Persistent,
	})

	if err != nil {
		logger.WithError(err).Warn("Failed to publish chat message to RabbitMQ")
	} else {
		logger.Info("Successfully published chat message to RabbitMQ")
	}
}

// FetchOfflineOnConnect retrieves and sends offline messages to a connected user
func FetchOfflineOnConnect(userId string, client *model.Client) {
	loggger := logutil.GetLogger()

	queue, err := MessageForOfflineUser.QueueDeclare(userId, true, false, false, false, nil)
	if err != nil {
		loggger.WithError(err).Warn("Failed to declare queue for offline user")
		return
	}

	err = MessageForOfflineUser.QueueBind(queue.Name, userId, "chat_exchange", false, nil)
	if err != nil {
		loggger.WithError(err).Warn("Failed to bind queue for offline user")
		return
	}

	queue, err = MessageForOfflineUser.QueueInspect(userId)
	if err != nil {
		loggger.WithError(err).Warn("Failed to inspect queue")
		return
	}

	if queue.Messages > 0 {
		consume, err := MessageForOfflineUser.Consume(
			userId,
			"",
			true,
			false,
			false,
			true,
			nil,
		)
		if err != nil {
			loggger.WithError(err).Warn("Failed to start consuming messages")
			return
		}

		for msg := range consume {
			var chat model.ChatMessage
			if err := json.Unmarshal(msg.Body, &chat); err == nil {
				if client.Conn != nil {
					err := client.Conn.WriteJSON(chat)
					if err != nil {
						loggger.WithError(err).Warn("Failed to send message to client")
					}
				} else {
					loggger.Warn("Cannot send message: client connection is nil")
				}
			} else {
				loggger.Warnf("Failed to unmarshal offline message: %s", string(msg.Body))
			}
		}
	}
}
