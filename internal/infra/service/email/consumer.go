package emailservice

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
	"gopkg.in/gomail.v2"
)

// RunConsumer starts the email consumption loop
func (s *Service) RunConsumer() error {
	fmt.Println("Starting consumer Email service")

	if s.rabbitmq == nil {
		return fmt.Errorf("rabbitmq service not initialized")
	}

	// Consumer uses an empty string as routing key for email exchange by default
	// or we can define a specific one if needed. Let's use "email" as RK for consistency.
	msgs, err := s.rabbitmq.ConsumeMessages(rabbitmq.EMAIL_EX, "")
	if err != nil {
		return fmt.Errorf("error starting consumer: %v", err)
	}

	log.Println("Email Worker: started and waiting for messages...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Email Worker: received a message: %s", d.Body)

			var bodyEmail BodyEmail
			if err := json.Unmarshal(d.Body, &bodyEmail); err != nil {
				log.Printf("Email Worker: error unmarshaling message: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := s.processSendEmail(&bodyEmail); err != nil {
				log.Printf("Email Worker: error sending email: %v", err)
				d.Nack(false, true) // Requeue if it fails to send
				continue
			}

			log.Printf("Email Worker: email sent successfully to %s", bodyEmail.Email)
			d.Ack(false)
		}
	}()

	<-forever
	return nil
}

func (s *Service) processSendEmail(bodyEmail *BodyEmail) error {
	// Configurações do e-mail via ENV (mesma lógica anterior)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("EMAIL_SERVICE")
	senderPass := os.Getenv("PASSWORD_EMAIL_SERVICE")

	if smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP configuration missing in environment variables")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("error converting email port: %v", err)
	}

	// Criar mensagem
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", bodyEmail.Email)
	m.SetHeader("Subject", bodyEmail.Subject)
	m.SetBody("text/html", bodyEmail.Body)

	// Configurar servidor SMTP
	d := gomail.NewDialer(smtpHost, port, senderEmail, senderPass)

	// Enviar o e-mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
