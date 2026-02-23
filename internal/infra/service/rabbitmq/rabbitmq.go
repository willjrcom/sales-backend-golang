package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	GROUP_ITEM_EX     = "print.group.item"
	ORDER_EX          = "print.order"
	ORDER_DELIVERY_EX = "print.order.delivery"
	EMAIL_EX          = "email"
)

// RabbitMQ structure to manage connection and channel
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

type PrintMessage struct {
	Id          string `json:"id"`
	PrinterName string `json:"printer_name"`
}

// NewInstance creates and returns a new instance of RabbitMQ with retries
func NewInstance(url string) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	// Try to connect and create a channel with retries
	for i := range 5 {
		conn, err = amqp.Dial(url)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				break
			}
			conn.Close()
		}
		log.Printf("Retrying RabbitMQ connection (attempt %d/5)...", i+1)
		time.Sleep(2 * time.Second)
	}

	// If we failed to connect after retries
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		url:     url,
	}, nil
}

// EnsureExchangeQueueAndBind ensures the exchange, queue, and binding exist
func (r *RabbitMQ) EnsureExchangeQueueAndBind(exchange, routingKey string) error {
	// Names for Exchange and Queue
	exchangeName := fmt.Sprintf("%s_exchange", exchange)          // Example: empresa_123_exchange
	queueName := fmt.Sprintf("%s_%s_queue", exchange, routingKey) // Example: empresa_123_impressao.pedido_queue

	// Declare the Exchange (direct type)
	err := r.channel.ExchangeDeclare(
		exchangeName, // Name of the exchange
		"direct",     // Type of exchange
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s", err)
	}

	// Declare the Queue (durable)
	_, err = r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	// Bind the Queue to the Exchange with the Routing Key
	err = r.channel.QueueBind(
		queueName,    // Queue name
		routingKey,   // Routing Key (topic)
		exchangeName, // Exchange name
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %s", err)
	}

	return nil
}

func (r *RabbitMQ) SendPrintMessage(exchage, routingKey, message, printerName string) error {
	printMessage := PrintMessage{
		Id:          message,
		PrinterName: printerName,
	}

	printMessageJSON, err := json.Marshal(printMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal print message: %s", err)
	}

	return r.SendMessage(exchage, routingKey, string(printMessageJSON))
}

// SendMessage sends a message to a specific company's exchange with a routing key
func (r *RabbitMQ) SendMessage(exchange, routingKey, message string) error {
	// Ensure the exchange, queue, and binding are created
	err := r.EnsureExchangeQueueAndBind(exchange, routingKey)
	if err != nil {
		return fmt.Errorf("failed to ensure exchange, queue and binding: %s", err)
	}

	// Publish the message to the exchange with the routing key
	exchangeName := fmt.Sprintf("%s_exchange", exchange)
	err = r.channel.Publish(
		exchangeName, // Exchange name
		routingKey,   // Routing key (topic)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}

	fmt.Println("Message sent to exchange:", exchangeName)
	return nil
}

// ConsumeMessages starts consuming messages from the specified company's queue
func (r *RabbitMQ) ConsumeMessages(exchange, routingKey string) (<-chan amqp.Delivery, error) {
	// Ensure the exchange, queue, and binding exist
	err := r.EnsureExchangeQueueAndBind(exchange, routingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure exchange, queue and binding: %s", err)
	}

	// Start consuming messages from the queue
	queueName := fmt.Sprintf("%s_%s_queue", exchange, routingKey) // Example: empresa_123_impressao.pedido_queue
	msgs, err := r.channel.Consume(
		queueName, // Queue name
		"",        // Consumer name
		false,     // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %s", err)
	}

	return msgs, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		log.Printf("Error closing channel: %s", err)
	}
	if err := r.conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	}
}
