package kafka

import (
	"log"
	"os"
	"os/signal"
	"testing"
)

type KafkaHandlerTest struct {
}

func (h *KafkaHandlerTest) HandlerFunc(msg *ConsumerMessage) {
	log.Printf("Receive msg, key: %v, value: %v, topic: %v \n", string(msg.Key), string(msg.Value), string(msg.Topic))
}

func (h *KafkaHandlerTest) Close() {

}

func TestConsumerExample(t *testing.T) {
	config := ConsumerConfig{
		Topic:           []string{"metax.account_status.event"},
		SeedBrokers:     []string{"localhost:9092"},
		ConsumerGroupID: "cg-id-1",
	}
	consumerGroup, err := NewKafkaConsumer(config)
	if err != nil {
		t.Fatal(err)
	}

	handlerTest := &KafkaHandlerTest{}

	consumerGroup.SetHandler(handlerTest)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	isExit := make(chan bool, 1)
	consumerGroup.Start(signals, isExit)
	<-isExit
}

func TestProducer(t *testing.T) {

}
