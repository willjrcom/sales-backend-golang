/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/willjrcom/sales-backend-go/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "sales-backend-go",
	Short: "A brief description of your application",
}

func main() {
	cmd.HttpserverCmd.PersistentFlags().StringP("port", "p", ":8080", "the port to connect to server")
	cmd.HttpserverCmd.PersistentFlags().StringP("environment", "e", "dev", "the environment to run the server in")
	rootCmd.AddCommand(cmd.HttpserverCmd)

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}

}
