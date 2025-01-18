/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	"github.com/willjrcom/sales-backend-go/internal/infra/modules"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

// httpserverCmd represents the httpserver command
var HttpserverCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("httpserver called")
		port, _ := cmd.Flags().GetString("port")
		environment, _ := cmd.Flags().GetString("environment")
		fmt.Println(environment)
		flag.Parse()
		chi := server.NewServerChi()

		s3Service := s3service.NewS3Client()
		cmd.Println("s3 loaded")

		// Load database
		db := database.NewPostgreSQLConnection()
		cmd.Println("db loaded")

		// Iniciar a goroutine para flush
		// go func() {
		// 	for {
		// 		// Esperar 1 segundo antes de chamar o Flush novamente
		// 		time.Sleep(10 * time.Second)
		// 		producerKafka.Flush(int(time.Second) * 15)
		// 	}
		// }()

		modules.MainModules(db, chi, s3Service)
		cmd.Println("modules loaded")

		if err := chi.StartServer(port); err != nil {
			panic(err)
		}
	},
}
