/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"flag"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	categoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/category_product"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	sizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size_category"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size_category"
)

// httpserverCmd represents the httpserver command
var HttpserverCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("httpserver called")
		port, _ := cmd.Flags().GetString("port")

		flag.Parse()
		ctx := context.Background()
		server := server.NewServerChi()

		// Load database
		db, err := database.NewPostgreSQLConnection(ctx)

		if err != nil {
			panic(err)
		}

		// Load repositories
		productRepo := productrepositorybun.NewProductRepositoryBun(db)
		categoryRepo := categoryrepositorybun.NewCategoryProductRepositoryBun(ctx, db)
		sizeRepo := sizerepositorybun.NewSizeCategoryRepositoryBun(ctx, db)

		// Load services
		productService := productusecases.NewService(productRepo, categoryRepo)
		categoryProductService := categoryproductusecases.NewService(categoryRepo)
		sizeService := sizeusecases.NewService(sizeRepo)

		// Load handlers
		productHandler := handlerimpl.NewHandlerProduct(productService)
		categoryHandler := handlerimpl.NewHandlerCategoryProduct(categoryProductService)
		sizeHandler := handlerimpl.NewHandlerSizeProduct(sizeService)

		server.AddHandler(productHandler)
		server.AddHandler(categoryHandler)
		server.AddHandler(sizeHandler)

		if err := server.StartServer(port); err != nil {
			panic(err)
		}
	},
}

func HandlerRegisterProduct(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("register product"))
}
