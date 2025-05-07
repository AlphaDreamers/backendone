package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/nats-io/nats.go"
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
	viper.AddConfigPath("/etc/myapp/")   // You can also add other directories like global config paths
	viper.AddConfigPath("$HOME/.myapp/") // Add home directory for user-specific configurations

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

func natsProvider(log *logrus.Logger, v *viper.Viper) *nats.Conn {
	natsUrl := v.GetString("nats.url")
	if natsUrl == "" {
		natsUrl = "nats://127.0.0.1:4222"
	}
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatalf("Error connecting to NATS, %s", err.Error())
	}
	return nc
}
func cleanUpDatabase(log *logrus.Logger, db *gorm.DB) {
	if db == nil {
		log.Warn("Database is nil, nothing to clean up")
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Failed to get sql.DB from GORM: %v", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Errorf("Failed to close database connection: %v", err)
	} else {
		log.Info("Database connection closed successfully")
	}
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
		AutoMigrate()
	if err != nil {
		log.Debugln("Failed to migrate DB")
		return err
	}
	return nil
}
