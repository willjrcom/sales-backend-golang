/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/willjrcom/sales-backend-go/internal/infra/pb/orderpb"
	"google.golang.org/grpc"
)

// grpcclientCmd represents the grpcclient command
var GrpcclientCmd = &cobra.Command{
	Use:   "grpcclient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Println("httpclient called")
		name, _ := cmd.Flags().GetString("name")

		// Set up a connection to the server.
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()
		c := orderpb.NewOrderServiceClient(conn)

		//r, err := c.GetAllOrder(ctx, &orderpb.BlankMessage{})
		r, err := c.Hello(context.Background(), &orderpb.Message{Name: name})

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetName())
	},
}
