package emailservice

import (
	"encoding/json"
	"fmt"

	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

type BodyEmail struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Service struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewService(rabbitmq *rabbitmq.RabbitMQ) *Service {
	return &Service{rabbitmq: rabbitmq}
}

func (s *Service) SendEmail(bodyEmail *BodyEmail) error {
	if s.rabbitmq == nil {
		return fmt.Errorf("rabbitmq service not initialized")
	}

	bodyJSON, err := json.Marshal(bodyEmail)
	if err != nil {
		return fmt.Errorf("error marshaling email body: %v", err)
	}

	// Publish to RabbitMQ using the common EMAIL_EX exchange
	// We use an empty routing key since it's a direct exchange for emails
	if err := s.rabbitmq.SendMessage(rabbitmq.EMAIL_EX, "", string(bodyJSON)); err != nil {
		return fmt.Errorf("error publishing email to rabbitmq: %v", err)
	}

	fmt.Println("Email message published to RabbitMQ successfully")
	return nil
}
