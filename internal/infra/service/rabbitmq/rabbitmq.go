package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	SHIFT_EX      = "print.shift"
	GROUP_ITEM_EX = "print.group.item"
	ORDER_EX      = "print.order"
	EMAIL_EX      = "email"
)

const (
	SHIFT_PATH      = "/print-manager/shift/"
	GROUP_ITEM_PATH = "/print-manager/group-item/"
	ORDER_PATH      = "/print-manager/order/"
)

// RabbitMQ structure to manage connection and channel
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	mu      sync.Mutex
}

type PrintMessage struct {
	Path        string `json:"path"`
	PrinterName string `json:"printer_name"`
}

// NewInstance creates and returns a new instance of RabbitMQ with retries
func NewInstance(url string) (*RabbitMQ, error) {
	r := &RabbitMQ{url: url}
	if err := r.connect(); err != nil {
		return nil, err
	}
	return r, nil
}

// connect handles the actual connection and channel creation with retries
func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	for i := 0; i < 5; i++ {
		r.conn, err = amqp.Dial(r.url)
		if err == nil {
			r.channel, err = r.conn.Channel()
			if err == nil {
				log.Println("Successfully connected to RabbitMQ")
				return nil
			}
			r.conn.Close()
		}
		log.Printf("Retrying RabbitMQ connection (attempt %d/5)...", i+1)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to RabbitMQ after retries: %v", err)
}

// reconnectIfClosed checks if connection or channel is closed and reconnects if necessary
func (r *RabbitMQ) reconnectIfClosed() error {
	if r.conn == nil || r.conn.IsClosed() || r.channel == nil {
		log.Println("RabbitMQ connection closed, attempting to reconnect...")
		return r.connect()
	}
	return nil
}

// EnsureExchangeQueueAndBind ensures the exchange, queue, and binding exist
func (r *RabbitMQ) EnsureExchangeQueueAndBind(exchange, routingKey string) error {
	if err := r.reconnectIfClosed(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Names for Exchange and Queue
	exchangeName := fmt.Sprintf("%s_exchange", exchange)
	queueName := fmt.Sprintf("%s_%s_queue", exchange, routingKey)

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

func (r *RabbitMQ) SendPrintMessage(exchage, routingKey, path, printerName string) error {
	if printerName == "" {
		printerName = "default"
	}
	printMessage := PrintMessage{
		Path:        path,
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

	r.mu.Lock()
	defer r.mu.Unlock()

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

	r.mu.Lock()
	defer r.mu.Unlock()

	// Start consuming messages from the queue
	queueName := fmt.Sprintf("%s_%s_queue", exchange, routingKey)
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
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %s", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}
}
