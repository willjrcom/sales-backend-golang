package emailservice

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendEmail(destination *string) error {
	// Configurações do e-mail
	smtpHost := "smtp.gmail.com" // Por exemplo, para Gmail
	smtpPort := 587              // Porta do servidor SMTP
	senderEmail := "email"
	senderPass := "pass" // Ou utilize tokens de aplicativo para maior segurança

	// Criar mensagem
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", *destination)
	m.SetHeader("Subject", "Teste de envio com GoMail")
	m.SetBody("text/plain", "Este é um teste de envio de e-mail utilizando a biblioteca GoMail.")

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
