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
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	addressrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/address"
	categoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/category_product"
	clientrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/client"
	contactrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/contact"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	sizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size_category"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact_person"
	deliveryorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_order"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
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
		categoryRepo := categoryrepositorybun.NewCategoryProductRepositoryBun(db)
		sizeRepo := sizerepositorybun.NewSizeCategoryRepositoryBun(db)

		clientRepo := clientrepositorybun.NewClientRepositoryBun(db)
		contactRepo := contactrepositorybun.NewContactRepositoryBun(db)
		addressRepo := addressrepositorybun.NewAddressRepositoryBun(db)

		orderRepo := orderrepositorybun.NewOrderRepositoryBun(db)
		deliveryOrderRepo := orderrepositorybun.NewDeliveryOrderRepositoryBun(db)

		// Load services
		productService := productusecases.NewService(productRepo, categoryRepo)
		categoryProductService := categoryproductusecases.NewService(categoryRepo)
		sizeService := sizeusecases.NewService(sizeRepo)

		clientService := clientusecases.NewService(clientRepo, contactRepo)
		contactService := contactusecases.NewService(contactRepo)

		orderService := orderusecases.NewService(orderRepo)
		deliveryOrderService := deliveryorderusecases.NewService(deliveryOrderRepo, addressRepo, clientRepo, orderRepo)

		// Load handlers
		productHandler := handlerimpl.NewHandlerProduct(productService)
		categoryHandler := handlerimpl.NewHandlerCategoryProduct(categoryProductService)
		sizeHandler := handlerimpl.NewHandlerSizeProduct(sizeService)

		clientHandler := handlerimpl.NewHandlerClient(clientService)
		contactHandler := handlerimpl.NewHandlerContactPerson(contactService)

		orderHandler := handlerimpl.NewHandlerOrder(orderService)
		deliveryOrderHandler := handlerimpl.NewHandlerDeliveryOrder(deliveryOrderService)

		server.AddHandler(productHandler)
		server.AddHandler(categoryHandler)
		server.AddHandler(sizeHandler)
		server.AddHandler(clientHandler)
		server.AddHandler(contactHandler)
		server.AddHandler(orderHandler)
		server.AddHandler(deliveryOrderHandler)

		if err := server.StartServer(port); err != nil {
			panic(err)
		}
	},
}
