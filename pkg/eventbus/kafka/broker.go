package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"
)

// Message - message info
type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

// EventData - Event Data Info
type EventData struct {
	RequestId string      `json:"request_id"`
	EventName string      `json:"event_name"`
	Payload   interface{} `json:"payload"`
	CreateAt  time.Time   `json:"create_at"`
}

type DataTest struct {
	Author string `json:"author"`
}

func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}
