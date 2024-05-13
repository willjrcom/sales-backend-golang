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
	clientrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/client"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	contactrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/contact"
	employeerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/employee"
	groupitemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/group_item"
	itemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/item"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	orderprocessrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order_process"
	orderqueuerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order_queue"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	productcategoryprocessrulerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_process_rule"
	productcategoryproductrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_product"
	productcategoryquantityrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_quantity"
	productcategorysizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_size"
	schemarepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/schema"
	shiftrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/shift"
	tablerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/table"
	userrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/user"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	itemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/item"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	orderdeliveryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_delivery"
	orderpickupusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_pickup"
	orderprocessusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_process"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
	ordertableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_table"
	productcategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
	productcategoryprocessruleusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_process_rule"
	productcategoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_product"
	productcategoryquantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_quantity"
	productcategorysizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_size"
	shiftusecases "github.com/willjrcom/sales-backend-go/internal/usecases/shift"
	tableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
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
		server := server.NewServerChi()

		s3Service := s3service.NewS3Client()

		// Load database
		db, err := database.NewPostgreSQLConnection(ctx)

		if err != nil {
			panic(err)
		}

		// Load repositories
		productRepo := productcategoryproductrepositorybun.NewProductRepositoryBun(db)
		productCategoryRepo := productcategoryrepositorybun.NewProductCategoryRepositoryBun(db)
		sizeRepo := productcategorysizerepositorybun.NewSizeRepositoryBun(db)
		quantityRepo := productcategoryquantityrepositorybun.NewQuantityRepositoryBun(db)
		processRuleRepo := productcategoryprocessrulerepositorybun.NewProcessRuleRepositoryBun(db)

		clientRepo := clientrepositorybun.NewClientRepositoryBun(db)
		contactRepo := contactrepositorybun.NewContactRepositoryBun(ctx, db)
		addressRepo := addressrepositorybun.NewAddressRepositoryBun(db)

		orderRepo := orderrepositorybun.NewOrderRepositoryBun(db)
		orderDeliveryRepo := orderrepositorybun.NewOrderDeliveryRepositoryBun(db)
		orderPickupRepo := orderrepositorybun.NewOrderPickupRepositoryBun(db)
		orderTableRepo := orderrepositorybun.NewOrderTableRepositoryBun(db)
		itemRepo := itemrepositorybun.NewItemRepositoryBun(db)
		groupItemRepo := groupitemrepositorybun.NewGroupItemRepositoryBun(db)

		orderQueueRepo := orderqueuerepositorybun.NewQueueRepositoryBun(db)
		processRepo := orderprocessrepositorybun.NewProcessRepositoryBun(db)

		employeeRepo := employeerepositorybun.NewEmployeeRepositoryBun(db)
		tableRepo := tablerepositorybun.NewTableRepositoryBun(db)
		shiftRepo := shiftrepositorybun.NewShiftRepositoryBun(db)

		schemaRepo := schemarepositorybun.NewSchemaRepositoryBun(db)
		companyRepo := companyrepositorybun.NewCompanyRepositoryBun(db)
		userRepo := userrepositorybun.NewUserRepositoryBun(db)

		// Load services
		productService := productcategoryproductusecases.NewService(productRepo, productCategoryRepo, s3Service)
		productCategoryService := productcategoryusecases.NewService(productCategoryRepo)
		sizeService := productcategorysizeusecases.NewService(sizeRepo, productCategoryRepo)
		quantityService := productcategoryquantityusecases.NewService(quantityRepo, productCategoryRepo)
		processRuleService := productcategoryprocessruleusecases.NewService(processRuleRepo)

		clientService := clientusecases.NewService(clientRepo, contactRepo)
		employeeService := employeeusecases.NewService(employeeRepo, contactRepo)
		contactService := contactusecases.NewService(contactRepo)

		itemService := itemusecases.NewService(itemRepo, groupItemRepo, orderRepo, productRepo, quantityRepo)
		groupItemService := groupitemusecases.NewService(itemRepo, groupItemRepo, productRepo)

		orderQueueService := orderqueueusecases.NewService(orderQueueRepo, processRepo)
		processService := orderprocessusecases.NewService(processRepo, groupItemRepo, orderQueueService, processRuleRepo, groupItemService)

		orderService := orderusecases.NewService(orderRepo, shiftRepo, groupItemService, processService, processRuleRepo, orderQueueService)
		orderPickupService := orderpickupusecases.NewService(orderPickupRepo, orderService)
		orderDeliveryService := orderdeliveryusecases.NewService(orderDeliveryRepo, addressRepo, clientRepo, orderRepo, employeeRepo, orderService)
		orderTableService := ordertableusecases.NewService(orderTableRepo, tableRepo, orderService)

		tableService := tableusecases.NewService(tableRepo)
		shiftService := shiftusecases.NewService(shiftRepo)

		schemaService := schemaservice.NewService(schemaRepo)
		userService := userusecases.NewService(userRepo)
		companyService := companyusecases.NewService(companyRepo, addressRepo, *schemaService, userRepo, *userService)

		// Load handlers
		productHandler := handlerimpl.NewHandlerProduct(productService)
		productCategoryHandler := handlerimpl.NewHandlerProductCategory(productCategoryService)
		sizeHandler := handlerimpl.NewHandlerSize(sizeService, productCategoryHandler.Path)
		quantityHandler := handlerimpl.NewHandlerQuantity(quantityService, productCategoryHandler.Path)
		processRuleHandler := handlerimpl.NewHandlerProcessRuleCategory(processRuleService, productCategoryHandler.Path)

		clientHandler := handlerimpl.NewHandlerClient(clientService)
		employeeHandler := handlerimpl.NewHandlerEmployee(employeeService)
		contactHandler := handlerimpl.NewHandlerContactPerson(contactService)

		orderHandler := handlerimpl.NewHandlerOrder(orderService)
		orderPickupHandler := handlerimpl.NewHandlerOrderPickup(orderPickupService)
		orderDeliveryHandler := handlerimpl.NewHandlerOrderDelivery(orderDeliveryService)
		orderTableHandler := handlerimpl.NewHandlerOrderTable(orderTableService)
		processHandler := handlerimpl.NewHandlerProcess(processService)
		orderQueueHandler := handlerimpl.NewHandlerQueue(orderQueueService)
		itemHandler := handlerimpl.NewHandlerItem(itemService)
		groupItemHandler := handlerimpl.NewHandlerGroupItem(groupItemService)

		tableHandler := handlerimpl.NewHandlerTable(tableService)
		shiftHandler := handlerimpl.NewHandlerShift(shiftService)

		companyHandler := handlerimpl.NewHandlerCompany(companyService)
		userHandler := handlerimpl.NewHandlerUser(userService)

		server.AddHandler(productHandler)
		server.AddHandler(productCategoryHandler)
		server.AddHandler(sizeHandler)
		server.AddHandler(quantityHandler)
		server.AddHandler(processRuleHandler)

		server.AddHandler(clientHandler)
		server.AddHandler(employeeHandler)
		server.AddHandler(contactHandler)

		server.AddHandler(orderHandler)
		server.AddHandler(orderPickupHandler)
		server.AddHandler(orderDeliveryHandler)
		server.AddHandler(orderTableHandler)
		server.AddHandler(processHandler)
		server.AddHandler(orderQueueHandler)
		server.AddHandler(itemHandler)
		server.AddHandler(groupItemHandler)

		server.AddHandler(tableHandler)
		server.AddHandler(shiftHandler)

		server.AddHandler(companyHandler)
		server.AddHandler(userHandler)

		if err := server.StartServer(port); err != nil {
			panic(err)
		}
	},
}
