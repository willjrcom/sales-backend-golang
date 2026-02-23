package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	emailservice "github.com/willjrcom/sales-backend-go/internal/infra/service/email"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

// EmailworkerCmd represents the emailworker command
var EmailworkerCmd = &cobra.Command{
	Use:   "emailworker",
	Short: "Runs the background worker to send emails from RabbitMQ",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("Email Worker starting...")

		conn := os.Getenv("RABBITMQ_URL")
		if conn == "" {
			fmt.Printf("RABBITMQ_URL not set")
			return
		}

		rabbitmqService, err := rabbitmq.NewInstance(conn)
		if err != nil {
			log.Fatalf("Error creating RabbitMQ instance: %s", err)
		}
		defer rabbitmqService.Close()

		emailService := emailservice.NewService(rabbitmqService)

		if err := emailService.RunConsumer(); err != nil {
			log.Fatalf("Email Worker failed: %s", err)
		}
	},
}
