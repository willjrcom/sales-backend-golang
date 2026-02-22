package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	GROUP_ITEM_QUEUE = "group_item_queue"
	ORDER_QUEUE      = "order_queue"
)

// RabbitMQ structure to manage connection and channel
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// NewInstance creates and returns a new instance of RabbitMQ with retries
func NewInstance(url string, maxRetries int, retryInterval time.Duration) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	// Try to connect and create a channel with retries
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				break
			}
			conn.Close()
		}
		log.Printf("Retrying RabbitMQ connection (attempt %d/%d)...", i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	// If we failed to connect after maxRetries
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %s", maxRetries, err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		url:     url,
	}, nil
}

// SendMessage sends a message to a specific company's exchange with a routing key
func (r *RabbitMQ) SendMessage(schemaName, routingKey, message string) error {
	exchangeName := fmt.Sprintf("%s_exchange", schemaName) // Exchange para a empresa

	// Ensure the exchange exists (direct type exchange)
	err := r.channel.ExchangeDeclare(
		exchangeName, // Name of the exchange
		"direct",     // Type of exchange (direct)
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s", err)
	}

	// Publish the message to the exchange with the routing key
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
func (r *RabbitMQ) ConsumeMessages(schemaName, queueName, routingKey string) (<-chan amqp.Delivery, error) {
	exchangeName := fmt.Sprintf("%s_exchange", schemaName) // Exchange da empresa

	// Ensure the queue exists
	_, err := r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare the queue: %s", err)
	}

	// Bind the queue to the exchange with the routing key (topic)
	err = r.channel.QueueBind(
		queueName,    // Queue name
		routingKey,   // Routing key (topic)
		exchangeName, // Exchange name
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %s", err)
	}

	// Start consuming messages from the queue
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
