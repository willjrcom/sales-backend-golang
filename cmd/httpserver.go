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
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	contactrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/contact"
	employeerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/employee"
	groupitemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/group_item"
	itemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/item"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	processrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/process_category"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	quantityrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/quantity_category"
	schemarepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/schema"
	shiftrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/shift"
	sizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size_category"
	tablerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/table"
	userrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/user"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/schema"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact"
	deliveryorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_order"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	itemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/item"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	processusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process_category"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity_category"
	shiftusecases "github.com/willjrcom/sales-backend-go/internal/usecases/shift"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size_category"
	tableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table"
	tableorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table_order"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
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
		quantityRepo := quantityrepositorybun.NewQuantityCategoryRepositoryBun(db)
		processRepo := processrepositorybun.NewProcessCategoryRepositoryBun(db)

		clientRepo := clientrepositorybun.NewClientRepositoryBun(db)
		contactRepo := contactrepositorybun.NewContactRepositoryBun(ctx, db)
		addressRepo := addressrepositorybun.NewAddressRepositoryBun(db)

		orderRepo := orderrepositorybun.NewOrderRepositoryBun(db)
		deliveryOrderRepo := orderrepositorybun.NewDeliveryOrderRepositoryBun(db)
		tableOrderRepo := orderrepositorybun.NewTableOrderRepositoryBun(db)
		itemRepo := itemrepositorybun.NewItemRepositoryBun(db)
		groupItemRepo := groupitemrepositorybun.NewGroupItemRepositoryBun(db)

		employeeRepo := employeerepositorybun.NewEmployeeRepositoryBun(db)
		tableRepo := tablerepositorybun.NewTableRepositoryBun(db)
		shiftRepo := shiftrepositorybun.NewShiftRepositoryBun(db)

		schemaRepo := schemarepositorybun.NewSchemaRepositoryBun(db)
		companyRepo := companyrepositorybun.NewCompanyRepositoryBun(db)
		userRepo := userrepositorybun.NewUserRepositoryBun(db)

		// Load services
		productService := productusecases.NewService(productRepo, categoryRepo)
		categoryProductService := categoryproductusecases.NewService(categoryRepo)
		sizeService := sizeusecases.NewService(sizeRepo)
		quantityService := quantityusecases.NewService(quantityRepo)
		processService := processusecases.NewService(processRepo)

		clientService := clientusecases.NewService(clientRepo, contactRepo)
		employeeService := employeeusecases.NewService(employeeRepo, contactRepo)
		contactService := contactusecases.NewService(contactRepo)

		orderService := orderusecases.NewService(orderRepo, shiftRepo)
		deliveryOrderService := deliveryorderusecases.NewService(deliveryOrderRepo, addressRepo, clientRepo, orderRepo, employeeRepo)
		tableOrderService := tableorderusecases.NewService(tableOrderRepo, tableRepo)
		itemService := itemusecases.NewService(itemRepo, groupItemRepo, orderRepo, productRepo, quantityRepo)
		groupService := groupitemusecases.NewService(itemRepo, groupItemRepo)

		tableService := tableusecases.NewService(tableRepo)
		shiftService := shiftusecases.NewService(shiftRepo)

		schemaService := schemaservice.NewService(schemaRepo)
		companyService := companyusecases.NewService(companyRepo, addressRepo, *schemaService)
		userService := userusecases.NewService(userRepo)

		// Load handlers
		productHandler := handlerimpl.NewHandlerProduct(productService)
		categoryHandler := handlerimpl.NewHandlerCategoryProduct(categoryProductService)
		sizeHandler := handlerimpl.NewHandlerSizeCategory(sizeService)
		quantityHandler := handlerimpl.NewHandlerQuantityCategory(quantityService)
		processHandler := handlerimpl.NewHandlerProcessCategory(processService)

		clientHandler := handlerimpl.NewHandlerClient(clientService)
		employeeHandler := handlerimpl.NewHandlerEmployee(employeeService)
		contactHandler := handlerimpl.NewHandlerContactPerson(contactService)

		orderHandler := handlerimpl.NewHandlerOrder(orderService)
		deliveryOrderHandler := handlerimpl.NewHandlerDeliveryOrder(deliveryOrderService)
		tableOrderHandler := handlerimpl.NewHandlerTableOrder(tableOrderService)
		itemHandler := handlerimpl.NewHandlerItem(itemService)
		groupHandler := handlerimpl.NewHandlerGroupItem(groupService)

		tableHandler := handlerimpl.NewHandlerTable(tableService)
		shiftHandler := handlerimpl.NewHandlerShift(shiftService)

		companyHandler := handlerimpl.NewHandlerCompany(companyService)
		userHandler := handlerimpl.NewHandlerUser(userService)

		server.AddHandler(productHandler)
		server.AddHandler(categoryHandler)
		server.AddHandler(sizeHandler)
		server.AddHandler(quantityHandler)
		server.AddHandler(processHandler)

		server.AddHandler(clientHandler)
		server.AddHandler(employeeHandler)
		server.AddHandler(contactHandler)

		server.AddHandler(orderHandler)
		server.AddHandler(deliveryOrderHandler)
		server.AddHandler(tableOrderHandler)
		server.AddHandler(itemHandler)
		server.AddHandler(groupHandler)

		server.AddHandler(tableHandler)
		server.AddHandler(shiftHandler)

		server.AddHandler(companyHandler)
		server.AddHandler(userHandler)

		if err := server.StartServer(port); err != nil {
			panic(err)
		}
	},
}
