package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
)

type KafkaProducer struct {
	*kafka.Producer
}

func NewProducer() *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
		panic(err)
	}

	return &KafkaProducer{p}
}

func (k *KafkaProducer) NewMessage(ctx context.Context, topic string, message interface{}) error {
	messageJson, err := json.Marshal(message)
	if err != nil {
		return errors.New("Failed to marshal message: " + err.Error())
	}

	schema, err := database.GetSchema(ctx)
	if err != nil {
		return err
	}

	kMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(schema),
		Value:          messageJson,
	}

	if err := k.Produce(kMessage, nil); err != nil {
		return errors.New("Failed to produce message: " + err.Error())
	}

	return nil
}
