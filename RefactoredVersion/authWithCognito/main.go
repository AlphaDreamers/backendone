package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SwanHtetAungPhyo/authCognito/cmd"
	"github.com/SwanHtetAungPhyo/authCognito/cmd/middleware"
	authMo "github.com/SwanHtetAungPhyo/authCognito/internal/handler/auth"
	uh "github.com/SwanHtetAungPhyo/authCognito/internal/handler/user"
	authRMo "github.com/SwanHtetAungPhyo/authCognito/internal/repo/auth"
	ur "github.com/SwanHtetAungPhyo/authCognito/internal/repo/user"
	authSMo "github.com/SwanHtetAungPhyo/authCognito/internal/service/auth"
	us "github.com/SwanHtetAungPhyo/authCognito/internal/service/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/gofiber/fiber/v2"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var InitProvideModule = fx.Module("initProvideModule", fx.Provide(
	LoadConfig,
	NewLogger,
	ConnectToRDS,
	NewFiberApp,
	LoadAwsConfig,
	NewCognitoClient,
	NewTextraClient,
	NewRekognitionClient,
	NewS3Client,
))

func main() {
	app := fx.New(
		InitProvideModule,
		authRMo.RepoModule,
		authSMo.ServiceLayerModule,
		authMo.HandlerLayerModule,
		ur.UserRepoModule,
		us.UserServiceModule,
		uh.UserHandlerModule,
		fx.Provide(
			cmd.NewAppState,
		),
		fx.Invoke(
			middleware.InitJWKS,
			cmd.AppLifeCycle),
	)

	app.Run()
}

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err.Error()))
	}
	return v
}

func NewLogger() *logrus.Logger {
	log := logrus.New()
	logFile := lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	multiwriter := io.MultiWriter(os.Stdout, &logFile)
	log.SetOutput(multiwriter)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", f.File, f.Line)
		},
		PrettyPrint: true,
	})
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	return log
}
func ConnectToRDS(v *viper.Viper, log *logrus.Logger) *gorm.DB {
	var db *gorm.DB
	var err error
	var once sync.Once
	var rawDb *sql.DB
	dsn := v.GetString("aws.rds.local")
	once.Do(func() {
		for i := 0; i < 10; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err == nil {
				break
			}
			rawDb, err = db.DB()
			if err != nil {
				log.Error("Failed to connect to RDS")
				return
			}
			log.Debugln("Connected to RDS")

			rawDb.SetMaxIdleConns(100)
			rawDb.SetMaxOpenConns(100)
			rawDb.SetConnMaxLifetime(time.Hour)

		}
	})
	return db
}

func NewFiberApp(v *viper.Viper, log *logrus.Logger) *fiber.App {
	idleTimeout := v.GetDuration("fiber.idleTimeout")
	readTimeout := v.GetDuration("fiber.readTimeout")
	writeTimeout := v.GetDuration("fiber.writeTimeout")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: v.GetBool("fiber.disableStartupMessage"),
		Prefork:               v.GetBool("fiber.prefork"),
		CaseSensitive:         v.GetBool("fiber.caseSensitive"),
		StrictRouting:         v.GetBool("fiber.strictRouting"),
		ServerHeader:          v.GetString("fiber.serverHeader"),
		AppName:               v.GetString("fiber.appName"),
		IdleTimeout:           idleTimeout,
		ReadTimeout:           readTimeout,
		WriteTimeout:          writeTimeout,
		BodyLimit:             20 * 1024 * 1024,
	})

	return app
}

func LoadAwsConfig() *aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err.Error())
	}
	return &cfg
}

func NewCognitoClient(awscfg *aws.Config) *cognitoidentityprovider.Client {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awscfg)
	return cognitoClient
}

func NewTextraClient(awscfg *aws.Config) *textract.Client {
	return textract.NewFromConfig(*awscfg)
}
func NewRekognitionClient(awscfg *aws.Config) *rekognition.Client {
	return rekognition.NewFromConfig(*awscfg)
}

func NewS3Client(awscfg *aws.Config) *s3.Client {
	return s3.NewFromConfig(*awscfg)
}
