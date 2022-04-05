package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"time"
)

const (
	serviceCode = "usms"
)

// ProducerType type
type ProducerType string

// Producer type
const (
	SyncProducerType  ProducerType = "sync"
	AsyncProducerType ProducerType = "async"
)

var (
	NoTopicDefined = errors.New("no topic defined")
)

// KafkaMessage - message info
type KafkaMessage struct {
	Topic string
	Key   []byte
	Value []byte
}

// ProducerConfig config
type ProducerConfig struct {
	SeedBrokers      []string
	NumFlushMessages int
	TopicMap         map[string]string
	Ssl              bool
	SslClientKey     string
	SslClientCert    string
	SslServerCert    string
}

// Producer --
type Producer struct {
	producer sarama.AsyncProducer
	topicMap map[string]string
}

// NewProducer init Producer
func NewProducer(cf ProducerConfig) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Flush.Messages = cf.NumFlushMessages
	config.Producer.Flush.Frequency = 1 * time.Second
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false

	if cf.Ssl {
		config.Net.TLS.Enable = true
		tlsConfig, err := NewTLSConfig(cf.SslClientCert, cf.SslClientKey, cf.SslServerCert)
		if err != nil {
			return nil, err
		}
		config.Net.TLS.Config = tlsConfig
	}

	asyncProducer, err := sarama.NewAsyncProducer(cf.SeedBrokers, config)
	if err != nil {
		zap.S().Infof("kafka async %v", err)
		return nil, err
	}

	prd := &Producer{
		producer: asyncProducer,
		topicMap: cf.TopicMap,
	}

	return prd, nil
}

// SendMessage send message to topic
func (p *Producer) SendMessage(m *KafkaMessage) {
	msg := &sarama.ProducerMessage{
		Topic: m.Topic,
		Key:   sarama.ByteEncoder(m.Key),
		Value: sarama.ByteEncoder(m.Value),
	}

	p.producer.Input() <- msg
}

func (p *Producer) SendAbstractMessage(msg interface{}) error {
	msgStructName := p.getTypeOfMessage(msg)
	topic := p.topicMap[msgStructName]
	if topic == "" {
		return NoTopicDefined
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(serviceCode),
		Value: sarama.ByteEncoder(msgBytes),
	}

	p.producer.Input() <- kafkaMsg
	er := <-p.producer.Errors()
	return er
}

func (p *Producer) getTypeOfMessage(msg interface{}) string {
	if t := reflect.TypeOf(msg); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

// Close topic
func (p *Producer) Close() {
	var wg sync.WaitGroup
	p.producer.AsyncClose()

	wg.Add(2)
	go func() {
		for range p.producer.Successes() {
			fmt.Println("Unexpected message on Successes()")
		}
		wg.Done()
	}()
	go func() {
		for msg := range p.producer.Errors() {
			fmt.Println(msg.Err)
		}
		wg.Done()
	}()
	wg.Wait()
}
