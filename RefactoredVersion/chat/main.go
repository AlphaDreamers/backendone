package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/SwanHtetAungPhyo/chat-order/cmd"
	"github.com/SwanHtetAungPhyo/chat-order/internal/handler/placeOrder"
	"github.com/SwanHtetAungPhyo/chat-order/internal/handler/ws"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/SwanHtetAungPhyo/chat-order/internal/repository"
	"github.com/SwanHtetAungPhyo/chat-order/internal/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	app := fx.New(
		InitModule,
		ws.WsModule,
		repository.RepoModule,
		service.ServiceModule,
		placeOrder.OrdHandlerModule,
		cmd.AppStateModule,
		fx.Invoke(
			//StartMigration,
			cmd.RegisterLifeCycle,
		),
	)
	app.Run()
}

func SetLogger() *logrus.Logger {
	logger := logrus.New()
	multiWriter := io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:   "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
	)

	logger.SetOutput(multiWriter)
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: true,
	})
	return logger
}

func LoadViper(log *logrus.Logger) *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	return v
}

func SetUpGorm(v *viper.Viper) *gorm.DB {
	dsn := v.GetString("aws.rds.local")
	if dsn == "" {
		logrus.Fatal("No RDS DSN provided")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{SingularTable: false},
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("Failed to get DB instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func LoadAWSConfig() *aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logrus.Fatalf("Failed to load AWS config: %v", err)
	}
	return &cfg
}

func DynamoDb(cfg *aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(*cfg)
}

func NewFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: false,
		Prefork:               false,
		StrictRouting:         false,
		CaseSensitive:         true,
		AppName:               "fiber",
	})
}

func RedisClient(v *viper.Viper) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     v.GetString("redis.addr"),
		Password: v.GetString("redis.password"),
		DB:       v.GetInt("redis.db"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

func NewS3Client(cfg *aws.Config) *s3.Client {
	return s3.NewFromConfig(*cfg)
}

func StartMigration(db *gorm.DB, log *logrus.Logger) {
	seedData := struct {
		Badges     []model.Badge
		Skills     []model.Skill
		GigTags    []model.GigTag
		Categories []model.Category
	}{
		Badges: []model.Badge{
			{Label: "Top Seller", Icon: "medal", Color: "gold"},
			{Label: "Fast Delivery", Icon: "rocket", Color: "blue"},
			{Label: "5-Star Rating", Icon: "star", Color: "purple"},
		},
		Skills: []model.Skill{
			{Label: "Web Development"},
			{Label: "Graphic Design"},
			{Label: "Digital Marketing"},
		},
		GigTags: []model.GigTag{
			{Label: "Frontend"},
			{Label: "Logo Design"},
			{Label: "SEO"},
		},
		Categories: []model.Category{
			{Label: "Digital Services", Slug: "digital-services", IsActive: true, SortOrder: 1},
			{Label: "Creative Design", Slug: "creative-design", IsActive: true, SortOrder: 2},
			{Label: "Marketing", Slug: "marketing", IsActive: true, SortOrder: 3},
		},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range []interface{}{
			seedData.Badges,
			seedData.Skills,
			seedData.Categories,
			seedData.GigTags,
		} {
			if err := tx.Create(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Errorf("Migration failed: %v", err)
	} else {
		log.Info("Initial data seeded successfully")
	}
}

var InitModule = fx.Module("init_module",
	fx.Provide(
		SetLogger,
		LoadViper,
		LoadAWSConfig,
		DynamoDb,
		RedisClient,
		NewFiberApp,
		NewS3Client,
		SetUpGorm,
	),
)
