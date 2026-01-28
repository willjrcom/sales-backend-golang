package emailservice

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type BodyEmail struct {
	Email   string
	Subject string
	Body    string
}

func SendEmail(bodyEmail *BodyEmail) error {
	// Configurações do e-mail
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("EMAIL_SERVICE")
	senderPass := os.Getenv("PASSWORD_EMAIL_SERVICE")

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("Erro ao converter a porta do email: %v", err)
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
		fmt.Printf("Erro ao enviar o e-mail: %v\n", err)
		return err
	}

	fmt.Println("E-mail enviado com sucesso!")
	return nil
}
