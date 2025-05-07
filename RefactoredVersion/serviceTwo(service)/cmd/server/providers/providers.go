package providers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"path"
	"runtime"
	"sync"
	"time"
)

var ProviderModule = fx.Provide(
	logProvider,
	LoadConfig,
	dataBaseProvider,
	//mongoProvider,
	NewFiberApp,
)

func logProvider() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: true,
	})
	return log
}
func LoadConfig(log *logrus.Logger) *viper.Viper {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	//viper.AddConfigPath("/etc/myapp/")
	//viper.AddConfigPath("$HOME/.myapp/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config.yaml file, %s", err.Error())
	}
	return viper.GetViper()
}

func dataBaseProvider(log *logrus.Logger, v *viper.Viper) *gorm.DB {
	var once sync.Once
	var dbInst *gorm.DB
	var tempDb *sql.DB
	var err error
	dsn := v.GetString("database.dsn")
	once.Do(func() {
		for i := 0; i < 10; i++ {
			dbInst, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
			tempDb, err = dbInst.DB()
			if err != nil {
				log.Debugln(err.Error())
			}
			tempDb.SetMaxIdleConns(10)
			tempDb.SetMaxOpenConns(100)
			tempDb.SetConnMaxLifetime(time.Hour)
			tempDb.SetConnMaxIdleTime(time.Hour)
		}

	})
	return dbInst
}

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		EnablePrintRoutes:     true,
	})
	return app
}

//func mongoProvider() *mongo.Client {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	uri := viper.GetString("mongo.uri")
//	clientOpts := options.Client().ApplyURI(uri)
//
//	client, err := mongo.Connect(ctx, clientOpts)
//	if err != nil {
//		log.Fatalf("MongoDB connection error: %v", err)
//	}
//
//	if err := client.Ping(ctx, nil); err != nil {
//		log.Fatalf("MongoDB ping failed: %v", err)
//	}
//
//	log.Println("MongoDB connection established.")
//	return client
//}
