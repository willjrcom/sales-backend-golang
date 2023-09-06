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
	rootCmd.PersistentFlags().StringP("port", "p", ":8080", "the port to connect to server")
	rootCmd.AddCommand(cmd.HttpserverCmd)

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}

}
