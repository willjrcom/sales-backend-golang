package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type KafkaConsumer struct {
	*kafka.Consumer
}

func NewConsumer() *KafkaConsumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "sales-backend-go",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
		panic(err)
	}

	return &KafkaConsumer{c}
}

func (c *KafkaConsumer) ReadMessages(topic string) {
	c.Subscribe(topic, nil)

	run := true

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Recebido sinal %v, finalizando...\n", sig)
			run = false
		default:
			msg, err := c.ReadMessage(time.Second * 120)
			if err == nil {
				process := &orderprocessentity.OrderProcess{}
				json.Unmarshal(msg.Value, process)
				fmt.Println(process)
				fmt.Printf("Mensagem recebida: %s\n", string(msg.Value))
			} else if err.(kafka.Error).Code() == kafka.ErrPartitionEOF {
				// Reached the end of the partition
				continue
			} else {
				fmt.Printf("Erro ao ler mensagem: %v (%v)\n", err, msg)
				run = false
			}

		}
	}

	time.Sleep(time.Second * 5)
}
