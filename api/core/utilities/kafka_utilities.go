package utilities

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

const (
	KafkaURL = "kafka-docker_kafka_1:9092"
	Topic    = "NewTopic1"
)

var count = 0

func PostDataToKafka(message string) {

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("Key-%d", count)),
		Value: []byte(message),
	}
	err := Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	}
	count = count + 1

}

func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
