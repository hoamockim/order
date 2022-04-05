package kafka

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

// InitOffsetType define const
type InitOffsetType int64

// InitOffsetNewest const
const (
	InitOffsetNewest InitOffsetType = -1
	InitOffsetOldest InitOffsetType = -2
)

// ConsumerConfig config
type ConsumerConfig struct {
	SeedBrokers     []string
	ConsumerGroupID string
	Topic           []string
	InitialOffset   InitOffsetType
	Ssl             bool
	SslClientKey    string
	SslClientCert   string
	SslServerCert   string
}

// ConsumerMessage message struct
type ConsumerMessage struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
	Timestamp time.Time
}

// KafkaConsumerHandlerFunc handle ...
type KafkaConsumerHandlerFunc func(message *ConsumerMessage)

type consumerGroupHandler struct {
	handleFunc KafkaConsumerHandlerFunc
}

// Setup ...
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup - clean
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim - claim message
func (cg consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		cg.handleFunc(&ConsumerMessage{
			Key:       msg.Key,
			Value:     msg.Value,
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Timestamp: msg.Timestamp,
		})
		sess.MarkMessage(msg, "")
	}
	return nil
}

// ConsumerHandler interface handler
type ConsumerHandler interface {
	HandlerFunc(*ConsumerMessage)
	Close()
}

// Consumer init consumer
type Consumer struct {
	client        sarama.Client
	consumerGroup sarama.ConsumerGroup
	Topic         []string
	Handler       ConsumerHandler
	running       bool
}

// NewConsumer - init consumer
func NewConsumer(cf ConsumerConfig) (*Consumer, error) {
	config := sarama.NewConfig()
	//config.Version = sarama.V1_0_0_0
	config.Version = sarama.V2_3_0_0
	config.Consumer.Return.Errors = true

	switch cf.InitialOffset {
	case InitOffsetNewest:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case InitOffsetOldest:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if cf.Ssl {
		config.Net.TLS.Enable = true
		tlsConfig, err := NewTLSConfig(cf.SslClientCert, cf.SslClientKey, cf.SslServerCert)
		if err != nil {
			return nil, err
		}
		config.Net.TLS.Config = tlsConfig
	}

	kkClient, err := sarama.NewClient(cf.SeedBrokers, config)
	if err != nil {
		return nil, err
	}

	kkConsumerGroup, err := sarama.NewConsumerGroupFromClient(cf.ConsumerGroupID, kkClient)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		client:        kkClient,
		consumerGroup: kkConsumerGroup,
		Topic:         cf.Topic,
	}

	return consumer, nil
}

// SetHandler controller consumer logic handler
func (c *Consumer) SetHandler(fn ConsumerHandler) {
	c.Handler = fn
}

// Start consumer
func (c *Consumer) Start(signals chan os.Signal, isExit chan bool) {
	ctx := context.Background()
	for {
		select {
		case <-signals:
			isExit <- true
			c.Close()
			break
		default:
			handler := consumerGroupHandler{handleFunc: c.Handler.HandlerFunc}
			err := c.consumerGroup.Consume(ctx, c.Topic, handler)
			if err != nil {
				zap.S().Warnf("consumer err %v %v", c.Topic, err)
				if !c.running {
					fmt.Println("Consumer closed")
				} else {
					fmt.Println("Error: ", err)
				}
			}
		}
	}
}

// Close consumer
func (c *Consumer) Close() {
	c.running = false
	_ = c.consumerGroup.Close()
	c.Handler.Close()
}
