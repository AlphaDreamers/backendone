package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	emailVerificationSubject = "user.email.verification"
	publishTimeout           = 5 * time.Second
)

type CoreNATSProducer struct {
	conn *nats.Conn
	log  *logrus.Logger
}

func NewCoreNATSProducer(nc *nats.Conn, log *logrus.Logger) *CoreNATSProducer {
	return &CoreNATSProducer{
		conn: nc,
		log:  log,
	}
}

func (p *CoreNATSProducer) PublishEmailVerification(email, token string) error {
	msg := model.EmailVerification{
		To:      email,
		Code:    token,
		Type:    "email-verification",
		Message: "This email verification code will expire in 10 minutes",
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Simple publish with timeout
	err = p.conn.Publish(emailVerificationSubject, msgBytes)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	// Ensure the message is sent
	if err := p.conn.FlushTimeout(publishTimeout); err != nil {
		return fmt.Errorf("flush failed: %w", err)
	}

	p.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Email verification published via core NATS")

	return nil
}
func StartCoreNATSSubscriber(nc *nats.Conn, log *logrus.Logger, handler func(email string, token string)) (*nats.Subscription, error) {
	sub, err := nc.Subscribe(emailVerificationSubject, func(m *nats.Msg) {
		var msg model.EmailVerification
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Errorf("Failed to unmarshal message: %v", err)
			return
		}
		handler(msg.To, msg.Code)
	})

	if err != nil {
		return nil, fmt.Errorf("subscription failed: %w", err)
	}

	log.Info("Core NATS subscriber started")
	return sub, nil
}
