package server

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/services"
	"github.com/hashicorp/consul/api"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"math/big"
	"path"
	"runtime"
	"sync"
	"time"
)

var ProviderModule = fx.Provide(
	logProvider,
	LoadConfig,
	dataBaseProvider,
	natsProvider,
	NewRedisClient,
	repo.NewAuthRepo,
	services.NewAuthService,
	ConsulRegistration,
	CertificateProvider,
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
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/myapp/")
	viper.AddConfigPath("$HOME/.myapp/")

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
			dbInst, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
					TablePrefix:   v.GetString("database.table_prefix"),
				},
			})
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
func natsProvider(log *logrus.Logger, v *viper.Viper) (*nats.Conn, error) {
	natsUrl := "nats://dummy:dummy@nats:4222"

	opts := []nats.Option{
		nats.MaxReconnects(5),
		nats.ReconnectWait(2 * time.Second),
		nats.Timeout(10 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Errorf("Disconnected from NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info("Reconnected to NATS")
		}),
	}

	// Add retry logic
	var nc *nats.Conn
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		nc, err = nats.Connect(natsUrl, opts...)
		if err == nil {
			break
		}
		log.Warnf("Failed to connect to NATS (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS after %d attempts: %w", maxRetries, err)
	}

	log.Info("Successfully connected to NATS")
	return nc, nil
}

func CertificateProvider(log *logrus.Logger, v *viper.Viper) *tls.Certificate {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"ASAuth Inc."},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("Failed to create TLS certificate: %v", err)
	}

	return &tlsCert
}

func ConsulRegistration(log *logrus.Logger, v *viper.Viper) struct{} {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = v.GetString("consul.address")

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
		return struct{}{}
	}

	serviceID := v.GetString("consul.service_id")
	serviceName := v.GetString("consul.service_name")
	servicePort := v.GetInt("consul.service_port")
	serviceTags := v.GetStringSlice("consul.service_tags")

	registration := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		Port: servicePort,
		Tags: serviceTags,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("https://localhost:%d/health", servicePort),
			Interval: "10s",
			Timeout:  "3s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	} else {
		log.Infof("Service registered with Consul: ID=%s, Name=%s, Port=%d, Tags=%v",
			serviceID, serviceName, servicePort, serviceTags)
	}
	return struct{}{}
}

func Migration(log *logrus.Logger, db *gorm.DB) error {
	err := db.Migrator().
		AutoMigrate(
			//&model.Country{},
			&model.User{},
		)
	if err != nil {
		log.Debugln("Failed to migrate DB")
		return err
	}
	return nil
}
func NewRedisClient(log *logrus.Logger, v *viper.Viper) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     v.GetString("redis.addr"),
		Password: v.GetString("redis.password"),
		DB:       v.GetInt("redis.database"),
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return redisClient
}
