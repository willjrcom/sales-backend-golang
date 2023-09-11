/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"flag"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category-product"
	productrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
)

const (
	defaultName = "world"
)

var ()

// httpserverCmd represents the httpserver command
var HttpserverCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("httpserver called")
		port, _ := cmd.Flags().GetString("port")

		flag.Parse()

		server := server.NewServerChi()

		// Load database
		_, err := database.NewPostgreSQLConnection()

		if err != nil {
			panic(err)
		}

		// Load repositories
		productRepo := productrepositorylocal.NewProductRepositoryLocal()
		categoryRepo := categoryrepositorylocal.NewCategoryProductRepositoryLocal()

		// Load services
		productService := productusecases.NewService(productRepo, categoryRepo)

		// Load handlers
		handlerProduct := handlerimpl.NewHandlerProduct(productService)

		server.AddHandler(handlerProduct)

		if err := server.StartServer(port); err != nil {
			panic(err)
		}
	},
}

func HandlerRegisterProduct(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("register product"))
}
