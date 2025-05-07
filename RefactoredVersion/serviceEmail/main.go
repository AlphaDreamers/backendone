package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/serviceEmail/consumer"
	"github.com/natefinch/lumberjack"
	"github.com/nats-io/nats.go"
	"github.com/resend/resend-go/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"path"
	"runtime"
	"strings"
	"time"
)

func main() {
	app := fx.New(
		// Providers
		fx.Provide(
			Logger,
			LoadViperConfig,
			NatsJetStream, // Changed from NatsClient to provide JetStream directly
			ResendClient,
			TopicProvider,
		),

		// Consumer Module
		consumer.Module,

		// Invoke consumer to start it
		fx.Invoke(func(*consumer.Consumer) {}),

		// Options
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
	)

	app.Run()
}

func Logger() *logrus.Logger {
	lumber := lumberjack.Logger{
		Filename:   "./email.log",
		MaxSize:    10,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   true,
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: false,
	})
	logger.SetOutput(&lumber)
	return logger
}

func LoadViperConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	return v
}

func NatsJetStream(v *viper.Viper, log *logrus.Logger) (nats.JetStreamContext, error) {
	natsUrl := v.GetString("nats.url")
	if natsUrl == "" {
		natsUrl = "nats://nats:4222" // Docker service name
	}

	opts := []nats.Option{
		nats.MaxReconnects(5),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Errorf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("Reconnected to NATS at %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Fatal("NATS connection permanently closed")
		}),
	}

	// Connect with retry
	var nc *nats.Conn
	var err error
	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(natsUrl, opts...)
		if err == nil {
			break
		}
		log.Warnf("NATS connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS after retries: %w", err)
	}

	// Create JetStream context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Ensure email stream exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "EMAIL_EVENTS",
		Subjects:  []string{"email.>"},
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
		Storage:   nats.FileStorage,
	})
	if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
		return nil, fmt.Errorf("failed to create email stream: %w", err)
	}

	// Ensure consumer exists
	_, err = js.AddConsumer("EMAIL_EVENTS", &nats.ConsumerConfig{
		Durable:       "email-consumer",
		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    3,
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, fmt.Errorf("failed to create email consumer: %w", err)
	}

	log.Info("NATS JetStream initialized successfully")
	return js, nil
}

func ResendClient(v *viper.Viper) *resend.Client {
	apiKey := v.GetString("resend.api_key")
	if apiKey == "" {
		panic("Resend API key not configured")
	}
	return resend.NewClient(apiKey)
}

func TopicProvider(v *viper.Viper) string {
	topic := v.GetString("nats.topic")
	if topic == "" {
		return "email.verification" // Default topic
	}
	return topic
}
