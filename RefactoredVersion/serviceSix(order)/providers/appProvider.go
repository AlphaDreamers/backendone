package providers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	rc "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/natefinch/lumberjack"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

var StartUpProvider = fx.Module(
	"startup providers",
	fx.Provide(
		OnceProvider,
		SetUpLogger,
		NatsClientProvider,
		FiberAppProvider,
		MongoDbProvider,
	),
)

func OnceProvider() *sync.Once {
	once := &sync.Once{}
	return once
}

func SetUpLogger(once *sync.Once) *logrus.Logger {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	logFile := &lumberjack.Logger{
		Filename:   path.Join(pwd, "app.log"), // Changed from ".log" to "app.log"
		MaxSize:    500,                       // Max size of log file in MB
		MaxBackups: 3,                         // Max number of backups
		MaxAge:     28,                        // Max age of logs in days
		Compress:   true,
	}

	var logger *logrus.Logger
	once.Do(func() {
		logger = logrus.New()
		multiWriter := io.MultiWriter(logFile, os.Stdout)
		logger.SetOutput(multiWriter)
		logger.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
			},
			TimestampFormat: time.RFC3339Nano,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			PrettyPrint: true,
		})
		logger.SetReportCaller(true)
		logger.SetLevel(logrus.DebugLevel)
	})
	return logger
}

func NatsClientProvider() *nats.Conn {
	client, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS: ", err.Error())
	}
	return client
}

func MongoDbProvider(once *sync.Once) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	once.Do(func() {
		clientURI := "mongodb://localhost:27017" // Replace with your MongoDB connection string
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(clientURI))
		if err != nil {
			log.Fatal("Error connecting to MongoDB: ", err.Error())
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal("Error pinging MongoDB: ", err.Error())
		}
	})
	return client
}

func FiberAppProvider(once *sync.Once) *fiber.App {
	var app *fiber.App
	once.Do(func() {
		app = fiber.New(fiber.Config{
			DisableStartupMessage: true,
			Prefork:               false,
			CaseSensitive:         true,
			StrictRouting:         true,
			ServerHeader:          "Fiber",
			AppName:               "Order Service",
			Concurrency:           100,
			WriteTimeout:          60 * time.Second,
			ReadTimeout:           60 * time.Second,
			IdleTimeout:           60 * time.Second,
			WriteBufferSize:       5 * 1024 * 1024,
			ReadBufferSize:        5 * 1024 * 1024,
		})

		app.Use(rc.New()) // Recovery middleware
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET, POST, PUT, PATCH, DELETE",
			AllowHeaders: "Origin, Content-Type, Accept",
		}))
		app.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Hello World",
			})
		})

		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}))
	})
	return app
}

func ViperProvider() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config.yaml file, %s", err)
	}
	return v
}
