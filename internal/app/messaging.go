package app

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Messaging struct {
	Enabled bool
	Kafka   Kafka // Kafka configuration
}

type Kafka struct {
	Brokers  []string
	User     string
	Password string   // For consumer group
	SSL      KafkaSSL // Enable SSL

	Producer *kafka.Writer
	Consumer *kafka.Reader // Consumer is created per-topic when needed
}

type KafkaSSL struct {
	Enabled    bool   // Enable SSL
	CertFile   string // Path to the certificate file
	KeyFile    string // Path to the key file
	CAFile     string // Path to the CA file
	ServerName string // Server name for SSL verification
}

// prepare sets up Kafka producer. Consumer is created per-topic when needed.
func (m *Messaging) prepare() error {
	if !m.Enabled {
		return nil
	}

	dialer, err := m.Kafka.buildDialer()
	if err != nil {
		return err
	}

	m.Kafka.Producer, err = m.Kafka.newProducer(dialer)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	m.Kafka.Producer.AllowAutoTopicCreation = true // Allow auto topic creation

	m.Kafka.Consumer, err = m.Kafka.newConsumer(dialer)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return nil
}

// ------- PRODUCER -------.
func (k *Kafka) newProducer(dialer *kafka.Dialer) (*kafka.Writer, error) {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.Brokers,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	}), nil
}

// ------- CONSUMER -------.
func (k *Kafka) newConsumer(dialer *kafka.Dialer) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.Brokers,
		GroupID:     "loan-service-consumer", // Consumer group ID
		GroupTopics: constant.ConsumerGroupTopics,
		Dialer:      dialer,

		// nolint
		MinBytes: 10e3, // 10KB.
		// nolint
		MaxBytes: 10e6, // 10MB
	})

	return reader, nil
}

// buildDialer handles SASL/SSL if configured.
func (k *Kafka) buildDialer() (*kafka.Dialer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second, //nolint
		DualStack: true,
	}

	// If SASL (username/password) is used
	if k.User != "" && k.Password != "" {
		plainMech := plain.Mechanism{
			Username: k.User,
			Password: k.Password,
		}
		dialer.SASLMechanism = plainMech
	}

	// If SSL is enabled
	if k.SSL.Enabled {
		tlsConfig, err := k.newTLSConfig(k.SSL)
		if err != nil {
			return nil, err
		}
		dialer.TLS = tlsConfig
	}

	return dialer, nil
}

// newTLSConfig loads certificates for SSL.
func (k *Kafka) newTLSConfig(cfg KafkaSSL) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed loading client cert/key: %w", err)
	}

	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading CA file: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certs")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   cfg.ServerName,
	}, nil
}
