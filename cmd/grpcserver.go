/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"flag"
	"log"
	"net"

	"github.com/spf13/cobra"
	grpcimpl "github.com/willjrcom/sales-backend-go/bootstrap/grpc"
	"github.com/willjrcom/sales-backend-go/internal/infra/pb/orderpb"
	"google.golang.org/grpc"
)

// grpcserverCmd represents the grpcserver command
var GrpcserverCmd = &cobra.Command{
	Use:   "grpcserver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("grpcserver called")
		flag.Parse()
		lis, err := net.Listen("tcp", ":50051")

		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		orderpb.RegisterOrderServiceServer(s, &grpcimpl.Server{})
		log.Printf("server listening at %v", lis.Addr())

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
