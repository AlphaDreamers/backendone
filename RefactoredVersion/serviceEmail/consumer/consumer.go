package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/serviceEmail/emails"
	"github.com/nats-io/nats.go"
	"github.com/resend/resend-go/v2"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"sync"
	"time"
)

type ConsumedMailAction struct {
	To      string   `json:"to"`
	Subject string   `json:"subject,omitempty"`
	Message string   `json:"message"`
	Type    string   `json:"type"`
	Link    string   `json:"link"`
	OrderID string   `json:"order_id"`
	Items   []string `json:"items"`
}

var Pool = sync.Pool{
	New: func() interface{} {
		return &ConsumedMailAction{}
	},
}

type Consumer struct {
	log      *logrus.Logger
	js       nats.JetStreamContext
	resend   *resend.Client
	subTopic string
	wg       sync.WaitGroup
	sub      *nats.Subscription
	stopChan chan struct{}
}

func NewConsumer(
	lc fx.Lifecycle,
	log *logrus.Logger,
	js nats.JetStreamContext,
	subTopic string,
	resend *resend.Client,
) (*Consumer, error) {
	if js == nil {
		return nil, errors.New("JetStream context is required")
	}

	consumer := &Consumer{
		log:      log,
		js:       js,
		resend:   resend,
		subTopic: subTopic,
		stopChan: make(chan struct{}),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return consumer.Start()
		},
		OnStop: func(ctx context.Context) error {
			return consumer.Stop()
		},
	})

	return consumer, nil
}

func (c *Consumer) Start() error {
	// Configure consumer with backoff strategy
	backoff := []time.Duration{100 * time.Millisecond, 500 * time.Millisecond, 1 * time.Second}

	var sub *nats.Subscription
	var err error

	for _, delay := range backoff {
		sub, err = c.js.QueueSubscribe(
			c.subTopic,
			"email-workers",
			c.handleMessage,
			nats.Durable("email-consumer"),
			nats.AckWait(30*time.Second),
			nats.MaxDeliver(5),
			nats.DeliverNew(),
			nats.ManualAck(),
		)

		if err == nil {
			break
		}

		c.log.WithError(err).Warnf("Failed to subscribe, retrying in %v", delay)
		time.Sleep(delay)
	}

	if err != nil {
		return fmt.Errorf("failed to subscribe after retries: %w", err)
	}

	c.sub = sub
	c.log.Info("JetStream consumer started successfully")
	return nil
}

func (c *Consumer) handleMessage(msg *nats.Msg) {
	select {
	case <-c.stopChan:
		return
	default:
	}

	action := Pool.Get().(*ConsumedMailAction)
	defer Pool.Put(action)

	if err := json.Unmarshal(msg.Data, action); err != nil {
		c.log.WithError(err).Error("Failed to unmarshal message")
		if ackErr := msg.Nak(); ackErr != nil {
			c.log.WithError(ackErr).Error("Failed to NAK message")
		}
		return
	}

	if err := c.processMessage(action); err != nil {
		c.log.WithError(err).Error("Failed to process message")
		if ackErr := msg.Nak(); ackErr != nil {
			c.log.WithError(ackErr).Error("Failed to NAK message")
		}
		return
	}

	if err := msg.Ack(); err != nil {
		c.log.WithError(err).Error("Failed to ACK message")
	}
}

func (c *Consumer) processMessage(action *ConsumedMailAction) error {
	c.wg.Add(1)
	defer c.wg.Done()

	var mail *resend.SendEmailRequest
	var err error

	switch action.Type {
	case "email-verification":
		if action.Link == "" {
			return errors.New("missing verification link")
		}
		mail = emails.VerificationEmail(action.To, action.Link)
	case "password-reset", "forgot-password":
		if action.Link == "" {
			return errors.New("missing password reset link")
		}
		mail = emails.PasswordResetEmail(action.To, action.Link)
	case "order-confirmation":
		if action.OrderID == "" {
			return errors.New("missing order ID")
		}
		mail = emails.OrderConfirmationEmail(action.To, action.OrderID, action.Items)
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}

	// Add subject if provided
	if action.Subject != "" {
		mail.Subject = action.Subject
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.resend.Emails.SendWithContext(ctx, mail)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	c.log.WithFields(logrus.Fields{
		"type": action.Type,
		"to":   action.To,
	}).Info("Email sent successfully")
	return nil
}

func (c *Consumer) Stop() error {
	c.log.Info("Shutting down consumer...")
	close(c.stopChan)

	if c.sub != nil {
		if err := c.sub.Drain(); err != nil {
			c.log.WithError(err).Error("Failed to drain subscription")
		}
	}

	// Wait for pending operations with timeout
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.log.Info("All pending emails processed")
	case <-time.After(30 * time.Second):
		c.log.Warn("Timeout waiting for pending emails")
	}

	c.log.Info("Consumer stopped successfully")
	return nil
}

var Module = fx.Module("email_consumer",
	fx.Provide(NewConsumer),
)
