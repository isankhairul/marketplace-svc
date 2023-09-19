package queue

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"marketplace-svc/helper/config"
	"os"
)

type KafkaBase struct {
	Producer    *kafka.Producer
	Consumer    *kafka.Consumer
	PrefixTopic string
	Log         logger.Logger
}

type handlerMessageSubscriber func(*kafka.Message)

func NewKafkaConsumer(cfg config.Config, log logger.Logger, groupID string) (*KafkaBase, error) {
	// add prefix topic
	groupID = cfg.Kafka.PrefixTopic + "." + groupID

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": cfg.Kafka.AutoOffset,
	})

	return &KafkaBase{Consumer: consumer, Log: log, PrefixTopic: cfg.Kafka.PrefixTopic}, err
}

func NewKafkaProducer(cfg config.Config) (*KafkaBase, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.BootstrapServers})
	return &KafkaBase{Producer: producer, PrefixTopic: cfg.Kafka.PrefixTopic}, err
}

func (k *KafkaBase) Publish(topic string, message []byte) error {
	// Delivery report handler for produced messages
	go func() {
		for e := range k.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// add prefix topic
	topic = k.PrefixTopic + "." + topic

	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	// Wait for message deliveries before shutting down
	k.Producer.Flush(15 * 1000)

	return err
}

func (k *KafkaBase) Subscribe(topic string, handler handlerMessageSubscriber) {
	// add prefix topic
	topic = k.PrefixTopic + "." + topic

	err := k.Consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		k.Log.Error(err)
		return
	}

	defer func(Consumer *kafka.Consumer) {
		err := Consumer.Close()
		if err != nil {
			k.Log.Error(err)
			return
		}
	}(k.Consumer)

	// A signal handler or similar could be used to set this to false to break the loop.
	var run = true
	for run {
		ev := k.Consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			_, err = k.Consumer.CommitMessage(e)
			if err == nil {
				handler(ev.(*kafka.Message))
			}
		case kafka.PartitionEOF:
			fmt.Printf("Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			run = false
		default:
			//fmt.Printf("Ignored %v\n", e)
		}
	}
}
