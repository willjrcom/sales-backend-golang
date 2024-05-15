/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"flag"

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

		flag.Parse()
		ctx := context.WithValue(context.Background(), database.Environment("environment"), environment)
		chi := server.NewServerChi()

		s3Service := s3service.NewS3Client()

		// Load database
		db, err := database.NewPostgreSQLConnection(ctx)

		if err != nil {
			panic(err)
		}

		modules.MainModules(db, chi, s3Service)

		if err := chi.StartServer(port); err != nil {
			panic(err)
		}
	},
}
