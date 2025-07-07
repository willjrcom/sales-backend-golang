package emailservice

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type BodyEmail struct {
	Email   string
	Subject string
	Body    string
}

func SendEmail(bodyEmail *BodyEmail) error {
	// Configurações do e-mail
	smtpHost := "smtp.gmail.com" // Por exemplo, para Gmail
	smtpPort := 587              // Porta do servidor SMTP
	senderEmail := os.Getenv("EMAIL_SERVICE")
	senderPass := os.Getenv("PASSWORD_EMAIL_SERVICE")

	// Criar mensagem
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", bodyEmail.Email)
	m.SetHeader("Subject", bodyEmail.Subject)
	m.SetBody("text/html", bodyEmail.Body)

	// Configurar servidor SMTP
	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPass)

	// Enviar o e-mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Erro ao enviar o e-mail: %v\n", err)
		return err
	}

	fmt.Println("E-mail enviado com sucesso!")
	return nil
}
