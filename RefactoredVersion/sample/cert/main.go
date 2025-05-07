package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CertificateProvider(log *logrus.Logger, v *viper.Viper) *tls.Certificate {
	// Generate ECDSA private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Certificate template
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

	// Create certificate from the template
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Marshal private key
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	// Define file paths for certificate and private key
	certFile := "certificates/cert.pem"
	keyFile := "certificates/key.pem"

	// Create the directory if it doesn't exist
	err = os.MkdirAll("certificates", 0755)
	if err != nil {
		log.Fatalf("Failed to create certificates directory: %v", err)
	}

	// Save the certificate to a file
	err = os.WriteFile(certFile, certPEM, 0644)
	if err != nil {
		log.Fatalf("Failed to write certificate to file: %v", err)
	}

	// Save the private key to a file
	err = os.WriteFile(keyFile, keyPEM, 0644)
	if err != nil {
		log.Fatalf("Failed to write private key to file: %v", err)
	}

	// Log paths of the saved files
	log.Infof("Certificate and private key saved to: %s, %s", certFile, keyFile)

	// Load and return the TLS certificate
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("Failed to create TLS certificate: %v", err)
	}

	return &tlsCert
}

func main() {
	// Example logrus logger and viper configuration
	log := logrus.New()
	v := viper.New()

	// Call the CertificateProvider to generate and save the certificate
	CertificateProvider(log, v)

	// You can now use the certificate in your server configuration
}
